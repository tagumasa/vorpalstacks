package broker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"

	"vorpalstacks/internal/services/aws/iot/ca"
)

type MessageHandler func(clientID string, topic string, payload []byte)

type Broker struct {
	server       *mqtt.Server
	handler      MessageHandler
	ca           *ca.CertificateAuthority
	authProvider AuthProvider
	running      atomic.Bool
	ctx          context.Context
	cancel       context.CancelFunc
	closeOnce    sync.Once
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

// SetAuthProvider configures the broker to use certificate and policy
// authentication. Must be called before Start.
func (b *Broker) SetAuthProvider(provider AuthProvider) {
	b.authProvider = provider
}

// SetMessageHandler configures the handler invoked for every published
// MQTT message. Must be called before Start.
func (b *Broker) SetMessageHandler(handler MessageHandler) {
	b.handler = handler
}

func (b *Broker) Start(port int) error {
	if !b.running.CompareAndSwap(false, true) {
		return fmt.Errorf("broker already running")
	}

	b.server = mqtt.New(&mqtt.Options{
		InlineClient: true,
	})

	if b.authProvider != nil {
		_ = b.server.AddHook(NewAuthHook(b.authProvider), nil)
	} else {
		// allowAllHook permits all connections when no AuthProvider is configured.
		allowAll := &allowAllHook{}
		_ = b.server.AddHook(allowAll, nil)
	}
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
				slog.Error("IoT MQTT broker panic recovered", "panic", r)
			}
		}()
		if err := b.server.Serve(); err != nil && b.running.Load() {
			slog.Error("IoT MQTT broker stopped", "error", err)
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

// allowAllHook permits all MQTT connections (used when no AuthProvider is set).
type allowAllHook struct {
	mqtt.HookBase
}

func (h *allowAllHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnACLCheck
}

func (h *allowAllHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	return true
}

func (h *allowAllHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	return true
}
