package kafka

import "context"

// KafkaProducerMock is a mock implementation of the Producer interface for testing
type KafkaProducerMock struct {
	PublishedMessages []string
}

// Publish simulates publishing a message to Kafka by appending it to PublishedMessages
func (kp *KafkaProducerMock) Publish(ctx context.Context, key, message string) error {
	// Simulate a message publish by storing it in the PublishedMessages slice
	kp.PublishedMessages = append(kp.PublishedMessages, message)
	return nil
}

// Close simulates closing the Kafka connection
func (kp *KafkaProducerMock) Close() error {
	// No operation needed, but implementing to satisfy the Producer interface
	return nil
}
