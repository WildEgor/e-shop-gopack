package mongo

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfiguer interface {
	URI() string
	DBName() string
}

type MongoConnection struct {
	client *mongo.Client
	cfg    MongoConfiguer
}

func NewMongoConnection() *MongoConnection {
	return &MongoConnection{}
}

func (mc *MongoConnection) Config(cfg MongoConfiguer) {
	mc.cfg = cfg
}

func (mc *MongoConnection) Connect(ctx context.Context) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mc.cfg.URI()))
	if err != nil {
		slog.Error("fail connect to mongo", err)
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		slog.Error("fail connect to mongo", err)
		panic(err)
	}

	slog.Info("success connect to mongoDB")

	mc.client = client
}

func (mc *MongoConnection) Disconnect(ctx context.Context) {
	if mc.client == nil {
		return
	}

	if err := mc.client.Disconnect(ctx); err != nil {
		slog.Error("fail disconnect to mongo", err)
		panic(err)
	}

	slog.Info("connection to mongo closed success")
}

func (mc *MongoConnection) DB() *mongo.Database {
	return mc.client.Database(mc.cfg.DBName())
}
