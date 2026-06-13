package broker

import (
	"fmt"
	"sync"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	"vorpalstacks/internal/services/aws/iot/ca"
)

type MessageHandler func(clientID string, topic string, payload []byte)

type Broker struct {
	server  *mqtt.Server
	handler MessageHandler
	ca      *ca.CertificateAuthority
	mu      sync.RWMutex
	running bool
}

func NewBroker(certificateAuthority *ca.CertificateAuthority, handler MessageHandler) *Broker {
	return &Broker{
		ca:      certificateAuthority,
		handler: handler,
	}
}

func (b *Broker) Start(port int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
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
		return fmt.Errorf("failed to add MQTT listener on port %d: %w", port, err)
	}

	go func() {
		if err := b.server.Serve(); err != nil {
			fmt.Printf("IoT MQTT broker stopped: %v\n", err)
		}
	}()

	b.running = true
	return nil
}

func (b *Broker) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running || b.server == nil {
		return nil
	}

	b.running = false
	return b.server.Close()
}

func (b *Broker) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

func (b *Broker) Publish(topic string, payload []byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if !b.running || b.server == nil {
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
