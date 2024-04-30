package notifier

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/wagslane/go-rabbitmq"
)

// TODO

type MockPublisher struct {
	mock.Mock
}

// Publish mocks the Publish method of *rabbitmq.Channel.
func (m *MockPublisher) PublishWithContext(
	ctx context.Context,
	data []byte,
	routingKeys []string,
	optionFuncs ...func(options *rabbitmq.PublisherOptions)) error {
	args := m.Called(ctx, data, routingKeys, optionFuncs)
	return args.Error(0)
}

// Close mocks the Close method of *rabbitmq.Channel.
func (m *MockPublisher) Close() error {
	args := m.Called()
	return args.Error(0)
}

var Config = &NotifierConfig{
	DSN:      "amqp://guest:guest@localhost:5672/",
	Exchange: "test",
}
