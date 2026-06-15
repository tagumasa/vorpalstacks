package broker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	"vorpalstacks/internal/services/aws/iot/ca"
)

type MessageHandler func(clientID string, topic string, payload []byte)

type Broker struct {
	server    *mqtt.Server
	handler   MessageHandler
	ca        *ca.CertificateAuthority
	running   atomic.Bool
	ctx       context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
}

func NewBroker(certificateAuthority *ca.CertificateAuthority, handler MessageHandler) *Broker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broker{
		ca:      certificateAuthority,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (b *Broker) Start(port int) error {
	if !b.running.CompareAndSwap(false, true) {
		return fmt.Errorf("broker already running")
	}

	b.server = mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	_ = b.server.AddHook(new(auth.AllowHook), nil)
	_ = b.server.AddHook(&iotHook{
		handler: b.handler,
	}, nil)

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "iot-tcp",
		Address: fmt.Sprintf(":%d", port),
	})

	if err := b.server.AddListener(tcp); err != nil {
		b.running.Store(false)
		return fmt.Errorf("failed to add MQTT listener on port %d: %w", port, err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("IoT MQTT broker panic recovered: %v\n", r)
			}
		}()
		if err := b.server.Serve(); err != nil && b.running.Load() {
			fmt.Printf("IoT MQTT broker stopped: %v\n", err)
		}
	}()

	return nil
}

func (b *Broker) Stop() error {
	if !b.running.CompareAndSwap(true, false) {
		return nil
	}

	b.cancel()

	b.closeOnce.Do(func() {
		if b.server != nil {
			_ = b.server.Close()
		}
	})

	return nil
}

func (b *Broker) IsRunning() bool {
	return b.running.Load()
}

func (b *Broker) Publish(topic string, payload []byte) error {
	if !b.running.Load() || b.server == nil {
		return fmt.Errorf("broker not running")
	}

	return b.server.Publish(topic, payload, false, 0)
}

func (b *Broker) Server() *mqtt.Server {
	return b.server
}

type iotHook struct {
	mqtt.HookBase
	handler MessageHandler
}

func (h *iotHook) ID() string { return "iot-hook" }

func (h *iotHook) Provides(b byte) bool {
	return b == mqtt.OnPublished
}

func (h *iotHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	if h.handler != nil {
		h.handler(cl.ID, pk.TopicName, pk.Payload)
	}
}
