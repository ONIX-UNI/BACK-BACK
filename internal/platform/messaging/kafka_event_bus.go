package messaging

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/DuvanRozoParra/sicou/pkg/messaging"
	"github.com/segmentio/kafka-go"
)

type KafkaEventBus struct {
	brokers []string
	writers map[string]*kafka.Writer
	mu      sync.Mutex
}

func NewKafkaEventBus(brokers []string) *KafkaEventBus {
	return &KafkaEventBus{
		brokers: brokers,
		writers: make(map[string]*kafka.Writer),
	}
}

func (k *KafkaEventBus) getWriter(topic string) *kafka.Writer {
	k.mu.Lock()
	defer k.mu.Unlock()

	if writer, exists := k.writers[topic]; exists {
		return writer
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(k.brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireAll,
	}

	k.writers[topic] = writer
	return writer
}

func (k *KafkaEventBus) Publish(ctx context.Context, event messaging.Event) error {

	topic := event.EventName()
	writer := k.getWriter(topic)

	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.EventKey()),
		Value: bytes,
	})
}

func (k *KafkaEventBus) Close() error {
	k.mu.Lock()
	defer k.mu.Unlock()

	for _, writer := range k.writers {
		_ = writer.Close()
	}

	return nil
}
