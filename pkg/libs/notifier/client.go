package notifier

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wagslane/go-rabbitmq"
)

type NotifierClient struct {
	topic string
	conn  *rabbitmq.Conn
	pub   *rabbitmq.Publisher
}

func NewNotifierClient(cfg *NotifierConfig) (*NotifierClient, error) {

	conn, err := rabbitmq.NewConn(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("new conn: %w", err)
	}

	pub, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(cfg.Exchange),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind("fanout"),
		rabbitmq.WithPublisherOptionsExchangeDurable,
	)
	if err != nil {
		return nil, fmt.Errorf("publisher: %w", err)
	}

	return &NotifierClient{
		topic: cfg.Exchange,
		conn:  conn,
		pub:   pub,
	}, nil
}

func (n *NotifierClient) Notify(ctx context.Context, payload *NotificationPayload) error {
	const contentTypeJSON = "application/json"

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return n.pub.PublishWithContext(
		ctx,
		body,
		[]string{"notifier.send-notification"},
		rabbitmq.WithPublishOptionsContentType(contentTypeJSON),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange(n.topic),
	)
}

func (n *NotifierClient) Close() error {
	n.pub.Close()

	if err := n.conn.Close(); err != nil {
		return fmt.Errorf("connection close: %w", err)
	}

	return nil
}
