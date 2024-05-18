package auth

import (
	"context"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	options *options

	Validate  func(ctx context.Context, in *pb.ValidateTokenRequest, opts ...grpc.CallOption) (*pb.UserData, error)
	FindByIds func(ctx context.Context, in *pb.FindByIdsRequest, opts ...grpc.CallOption) (*pb.FindByIdsResponse, error)
}

func WithDSN(dsn string) Option {
	return func(o *options) {
		o.DSN = dsn
	}
}

// NewClientWithOptions alternative
func NewClientWithOptions(opts ...Option) *Client {
	options := &options{}
	for _, o := range opts {
		o(options)
	}

	return &Client{
		options: options,
	}
}

// NewClient can be injectable via wire
func NewClient(cfg IOptions) *Client {
	options := &options{}
	opts := cfg.Options()
	options.DSN = opts.DSN

	return &Client{
		options: options,
	}
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(c.options.DSN, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return ErrEstablishGRPCConnect
	}

	service := pb.NewAuthServiceClient(conn)
	c.Validate = service.ValidateToken
	c.FindByIds = service.FindByIds

	return nil
}
