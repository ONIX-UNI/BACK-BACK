package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	Reader  *kafka.Reader
	brokers []string
}

func NewKafkaConsumer(
	brokers []string,
	topic string,
	groupID string,
) *KafkaConsumer {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &KafkaConsumer{
		Reader:  reader,
		brokers: brokers,
	}
}

func (c *KafkaConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.Reader.ReadMessage(ctx)
			if err != nil {
				return err
			}

			if err := handler(msg.Value); err != nil {
				return err
			}
		}
	}

}

func (c *KafkaConsumer) Close() error {
	return c.Reader.Close()
}

func (c *KafkaConsumer) HealthCheck(ctx context.Context) error {
	dialer := &kafka.Dialer{
		Timeout: 5 * time.Second,
	}

	conn, err := dialer.DialContext(ctx, "tcp", c.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Controller()
	return err
}

func (c *KafkaConsumer) EnsureTopic(
	ctx context.Context,
	topic string,
	partitions int,
	replicationFactor int,
) error {

	dialer := &kafka.Dialer{
		Timeout: 5 * time.Second,
	}

	// Conectarse a un broker
	conn, err := dialer.DialContext(ctx, "tcp", c.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	// Obtener controller
	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := dialer.DialContext(
		ctx,
		"tcp",
		fmt.Sprintf("%s:%d", controller.Host, controller.Port),
	)
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	// Intentar crear topic
	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	})

	// Si ya existe, no es error real
	if err != nil && err.Error() != "topic already exists" {
		return err
	}

	return nil
}
