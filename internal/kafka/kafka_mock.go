package kafka

import "context"

type KafkaProducerMock struct {
	PublishedMessages []string
}

func (kp *KafkaProducerMock) Publish(ctx context.Context, key, message string) error {

	kp.PublishedMessages = append(kp.PublishedMessages, message)
	return nil
}

func (kp *KafkaProducerMock) Close() error {

	return nil
}
