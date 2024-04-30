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

func NewPostgresConnection(
	cfg PostgresConfiguer,
) *PostgresConnection {

	conn := &PostgresConnection{
		cfg: cfg,
	}

	conn.Connect()

	return conn
}

// Connect make connect and ping db
func (p *PostgresConnection) Connect() {
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

	if err := dbpool.Ping(context.TODO()); err != nil {
		p.DB.Close()
	}
}

// Close close connection
func (p *PostgresConnection) Close() {
	if p.DB != nil {
		p.DB.Close()
	}

	slog.Info("connection to postgres closed success")
}
