package notifier_tests

import (
	"context"
	"github.com/WildEgor/g-core/pkg/libs/notifier"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateNotifier(t *testing.T) {
	assert.NotNil(t, client)
}

func TestNotifyNotifier(t *testing.T) {
	err := client.Notify(context.Background(), &notifier.NotificationPayload{})

	assert.Nil(t, err)
}
