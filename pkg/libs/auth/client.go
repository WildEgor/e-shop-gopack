package auth

import (
	"context"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
)

var _ IClient = (*Client)(nil)

type Client struct {
	options *options
	service pb.AuthServiceClient
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

	c.service = pb.NewAuthServiceClient(conn)

	return nil
}

func (c *Client) Validate(ctx context.Context, token string) (*pb.UserData, error) {
	res, err := c.service.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		st := status.Convert(err)
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			default:
				slog.Error("auth err", slog.Any("err", t))
				return nil, ErrUserNotFound
			}
		}
	}

	return res, nil
}

func (c *Client) FindUsersByIds(ctx context.Context, ids []string) ([]*pb.UserData, error) {
	res, err := c.service.FindByIds(ctx, &pb.FindByIdsRequest{
		Ids: ids,
	})

	if err != nil {
		st := status.Convert(err)
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			default:
				slog.Error("auth err", slog.Any("err", t))
				return []*pb.UserData{}, ErrUserNotFound
			}
		}
	}

	return res.Users, nil
}
