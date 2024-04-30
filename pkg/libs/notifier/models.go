package notifier

type EmailSettings struct {
	Email    string      `json:"email"`
	Subject  string      `json:"subject"`
	Text     string      `json:"text,omitempty"`
	Template string      `json:"template,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type PhoneSettings struct {
	Phone    string      `json:"phone"`
	Text     string      `json:"text"`
	Template string      `json:"template,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type PushSettings struct {
	To       string      `json:"to"`              // Must be unique identifier for each user, for example user_id
	Image    string      `json:"image,omitempty"` // Must be valid direct link to image
	Title    string      `json:"title,omitempty"`
	Message  string      `json:"message,omitempty"`
	Template string      `json:"template,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}

type NotificationPayload struct {
	Type          string        `json:"type"`
	Data          interface{}   `json:"data"`
	EmailSettings EmailSettings `json:"email_setting,omitempty"`
	PhoneSettings PhoneSettings `json:"phone_setting,omitempty"`
	PushSettings  PushSettings  `json:"push_settings,omitempty"`
}

type NotifierConfig struct {
	DSN      string `json:"dsn"`
	Exchange string `json:"exchange"`
}
