package notifier

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Email_Builder(t *testing.T) {
	notification := NewNotification(WithEmailSettings(EmailSettings{
		Subject: "subject",
		Text:    "text",
		Email:   "test@mail.ru",
	}))

	assert.Equal(t, "subject", notification.EmailSettings.Subject)
	assert.Equal(t, "text", notification.EmailSettings.Text)
	assert.Equal(t, "test@mail.ru", notification.EmailSettings.Email)
}

func Test_Phone_Builder(t *testing.T) {
	notification := NewNotification(WithPhoneSettings(PhoneSettings{
		Text:  "text",
		Phone: "phone",
	}))

	assert.Equal(t, "phone", notification.PhoneSettings.Phone)
	assert.Equal(t, "text", notification.PhoneSettings.Text)
}
