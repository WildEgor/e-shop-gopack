package auth

import "errors"

var (
	ErrEstablishGRPCConnect = errors.New("fail connect grpc")
)

type Option func(o *options)

type Options struct {
	DSN string
}

type options struct {
	Options
}
