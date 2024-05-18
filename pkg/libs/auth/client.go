package auth

import (
	"github.com/WildEgor/e-shop-gopack/pkg/libs/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	API     pb.AuthServiceClient
	options *options
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
func NewClient(cfg *Options) *Client {
	options := &options{}
	options.DSN = cfg.DSN

	return &Client{
		options: options,
	}
}

func (c *Client) Connect() error {
	conn, err := grpc.Dial(c.options.DSN, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return ErrEstablishGRPCConnect
	}

	c.API = pb.NewAuthServiceClient(conn)
	return nil
}
