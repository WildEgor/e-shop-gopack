package redis

import (
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

func NewRedisDBConnection(
	cfg RedisConfiguer,
) *RedisConnection {
	conn := &RedisConnection{
		nil,
		cfg,
	}

	conn.Connect()

	return conn
}

func (rc *RedisConnection) Connect() {
	opt, err := redis.ParseURL(rc.cfg.URI())
	if err != nil {
		slog.Error("fail parse url", err)
		panic(err)
	}

	rc.client = redis.NewClient(opt)

	if _, err := rc.client.Ping().Result(); err != nil {
		slog.Error("fail connect to redis ", err)
		panic(err)
	}

	slog.Info("success connect to redis")
}

func (rc *RedisConnection) Disconnect() {
	if rc.client == nil {
		return
	}

	if err := rc.client.Close(); err != nil {
		slog.Error("fail disconnect redis", err)
		return
	}

	slog.Info("connection to redis closed success")
}

func (rc *RedisConnection) Client() *redis.Client {
	return rc.client
}
