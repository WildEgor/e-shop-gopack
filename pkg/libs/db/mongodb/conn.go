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

func NewMongoConnection(
	cfg MongoConfiguer,
) *MongoConnection {
	conn := &MongoConnection{
		nil,
		cfg,
	}

	conn.Connect()

	return conn
}

func (mc *MongoConnection) Connect() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mc.cfg.URI()))
	if err != nil {
		slog.Error("fail connect to mongo", err)
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		slog.Error("fail connect to mongo", err)
		panic(err)
	}

	slog.Info("success connect to mongoDB")

	mc.client = client
}

func (mc *MongoConnection) Disconnect() {
	if mc.client == nil {
		return
	}

	if err := mc.client.Disconnect(context.TODO()); err != nil {
		slog.Error("fail disconnect to mongo", err)
		panic(err)
	}

	slog.Info("connection to mongo closed success")
}

func (mc *MongoConnection) DB() *mongo.Database {
	return mc.client.Database(mc.cfg.DBName())
}
