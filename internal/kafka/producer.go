package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, key, message string) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *KafkaProducer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{broker},
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	})
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) Publish(ctx context.Context, key, message string) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	})
	if err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
		return err
	}
	log.Printf("Successfully published message to Kafka - Key: %s, Message: %s", key, message)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
