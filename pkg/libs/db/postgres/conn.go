package postgres

import (
	"context"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
)

type PostgresConfiguer interface {
	URI() string
}

// PostgresConnection holds db conn
type PostgresConnection struct {
	DB  *pgxpool.Pool
	cfg PostgresConfiguer
}

func NewPostgresConnection() *PostgresConnection {
	return &PostgresConnection{}
}

func (p *PostgresConnection) Config(cfg PostgresConfiguer) {
	p.cfg = cfg
}

// Connect make connect and ping db
func (p *PostgresConnection) Connect(ctx context.Context) {
	config, err := pgxpool.ParseConfig(p.cfg.URI())
	if err != nil {
		slog.Error("fail parse config", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	dbpool, err := pgxpool.NewWithConfig(context.TODO(), config)
	if err != nil {
		slog.Error("fail connect to postgres", models.LogEntryAttr(&models.LogEntry{
			Err: err,
		}))
		panic(err)
	}

	p.DB = dbpool

	if err := dbpool.Ping(ctx); err != nil {
		p.DB.Close()
	}
}

// Close close connection
func (p *PostgresConnection) Disconnect(ctx context.Context) {
	if p.DB != nil {
		p.DB.Close()
	}

	slog.Info("connection to postgres closed success")
}
