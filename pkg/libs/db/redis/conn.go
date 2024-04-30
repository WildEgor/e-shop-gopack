package redis

import (
	"context"
	"github.com/go-redis/redis"
	"log/slog"
)

type RedisConfiguer interface {
	URI() string
}

type RedisConnection struct {
	client *redis.Client
	cfg    RedisConfiguer
}

func NewRedisDBConnection() *RedisConnection {
	return &RedisConnection{}
}

func (rc *RedisConnection) Config(cfg RedisConfiguer) {
	rc.cfg = cfg
}

func (rc *RedisConnection) Connect(ctx context.Context) {
	opt, err := redis.ParseURL(rc.cfg.URI())
	if err != nil {
		slog.Error("fail parse url", err)
		panic(err)
	}

	rc.client = redis.NewClient(opt)

	if _, err := rc.client.WithContext(ctx).Ping().Result(); err != nil {
		slog.Error("fail connect to redis ", err)
		panic(err)
	}

	slog.Info("success connect to redis")
}

func (rc *RedisConnection) Disconnect(ctx context.Context) {
	if rc.client == nil {
		return
	}

	if err := rc.client.WithContext(ctx).Close(); err != nil {
		slog.Error("fail disconnect redis", err)
		return
	}

	slog.Info("connection to redis closed success")
}

func (rc *RedisConnection) Client() *redis.Client {
	return rc.client
}
