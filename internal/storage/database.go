package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/logger"
	"strings"
)

type DB interface {
	Exec(ctx context.Context, query Query) (pgconn.CommandTag, error)
	Query(ctx context.Context, query Query) (pgx.Rows, error)
	QueryRow(ctx context.Context, query Query) pgx.Row
	Close() error
}

type Query struct {
	QueryName string
	Query     string
	Args      []any
}

func (q *Query) String() string {
	queryString := fmt.Sprintf("sql: %s: query: %s", q.QueryName, q.Query)
	if len(q.Args) != 0 {
		for k, v := range q.Args {
			queryString = strings.Replace(queryString, fmt.Sprintf("$%d", k+1), fmt.Sprintf("%v", v), 1)
		}
	}

	return queryString
}

type db struct {
	cfg     config.DBConfig
	connect *pgx.Conn
	log     logger.Logger
}

func (db *db) Close() error {
	return db.connect.Close(context.Background())
}

func (db *db) Exec(ctx context.Context, query Query) (pgconn.CommandTag, error) {
	db.logQuery(query)
	return db.connect.Exec(ctx, query.Query, query.Args...)
}
func (db *db) Query(ctx context.Context, query Query) (pgx.Rows, error) {
	db.logQuery(query)
	return db.connect.Query(ctx, query.Query, query.Args...)
}
func (db *db) QueryRow(ctx context.Context, query Query) pgx.Row {
	db.logQuery(query)
	return db.connect.QueryRow(ctx, query.Query, query.Args...)
}

func New(ctx context.Context, log logger.Logger, cfg config.DBConfig) (DB, error) {
	const op = "storage.New"
	connect, err := pgx.Connect(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("op: %s, error connect to database: %s", op, err.Error())
	}
	storage := &db{
		cfg:     cfg,
		connect: connect,
		log:     log,
	}

	if err := storage.connect.Ping(ctx); err != nil {
		return nil, fmt.Errorf("op: %s, error connect to database: %s", op, err.Error())
	}

	return storage, nil
}

func (db *db) logQuery(query Query) {
	db.log.Debug(query.String())
}
