package notifier_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/WildEgor/g-core/pkg/libs/notifier"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

var (
	cfg    *notifier.NotifierConfig
	client *notifier.NotifierClient
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	rmqc, err := rabbitmq.RunContainer(ctx,
		testcontainers.WithImage("rabbitmq:3-management-alpine"),
		rabbitmq.WithAdminUsername("admin"),
		rabbitmq.WithAdminPassword("root"),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	port, err := rmqc.MappedPort(ctx, "5672")
	if err != nil {
		log.Fatalf("Could not get RabbitMQ container port: %s", err)
	}

	dsn := fmt.Sprintf("amqp://guest:guest@127.0.0.1:%s/", port.Port())

	cfg = &notifier.NotifierConfig{
		DSN:      dsn,
		Exchange: "notifier",
	}

	defer func() {
		if err := rmqc.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}()

	client, err = notifier.NewNotifierClient(cfg)
	if err != nil {
		log.Fatalf("failed to create client container: %s", err)
	}
	defer func(client *notifier.NotifierClient) {
		err := client.Close()
		if err != nil {
			log.Fatalf("failed to create client container: %s", err)
		}
	}(client)

	code := m.Run()

	os.Exit(code)
}
