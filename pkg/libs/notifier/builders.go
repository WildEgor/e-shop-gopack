package notifier

type Option func(o *options)

type options struct {
	NotificationPayload
}

func WithData(data interface{}) Option {
	return func(o *options) { o.Data = data }
}

func WithEmailSettings(s EmailSettings) Option {
	return func(o *options) {
		o.Type = "email"
		o.EmailSettings = s
	}
}

func WithPhoneSettings(s PhoneSettings) Option {
	return func(o *options) {
		o.Type = "sms"
		o.PhoneSettings = s
	}
}

func NewNotification(opts ...Option) NotificationPayload {
	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return options.NotificationPayload
}
