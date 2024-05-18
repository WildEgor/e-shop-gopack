package auth

import (
	"context"
	"errors"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/auth/pb"
)

var (
	ErrEstablishGRPCConnect = errors.New("fail connect grpc")
	ErrUserNotFound         = errors.New("user not found")
	ErrUnknown              = errors.New("unknown error")
)

type Option func(o *options)

type Options struct {
	DSN string
}

type options struct {
	Options
}

type IOptions interface {
	Options() *Options
}

type IClient interface {
	Connect() error
	Validate(ctx context.Context, token string) (*pb.UserData, error)
	FindUsersByIds(ctx context.Context, ids []string) ([]*pb.UserData, error)
}
