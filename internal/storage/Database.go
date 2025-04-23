package storage

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/logger"
	"strings"
)

// todo add query in prarams
// todo add connect

type DB interface {
	ExecContext(ctx context.Context, query Query) (sql.Result, error)
	QueryContext(ctx context.Context, query Query) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query Query) *sql.Row
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
	cfg       config.DBConfig
	dbConnect *sql.DB
	log       logger.Logger
}

func (db *db) Close() error {
	return db.dbConnect.Close()
}

func (db *db) ExecContext(ctx context.Context, query Query) (sql.Result, error) {
	db.logQuery(query)
	return db.dbConnect.ExecContext(ctx, query.Query, query.Args...)
}
func (db *db) QueryContext(ctx context.Context, query Query) (*sql.Rows, error) {
	db.logQuery(query)
	return db.dbConnect.QueryContext(ctx, query.Query, query.Args...)
}
func (db *db) QueryRowContext(ctx context.Context, query Query) *sql.Row {
	db.logQuery(query)
	return db.dbConnect.QueryRowContext(ctx, query.Query, query.Args...)
}

func NewQuery(queryName string, query string, args []any) Query {
	return Query{
		QueryName: queryName,
		Query:     query,
		Args:      args,
	}
}

func NewDB(ctx context.Context, log logger.Logger, config config.DBConfig) DB {
	connect, err := sql.Open("sqlite3", config.DSN())
	if err != nil {
		log.Fatal("error connect to database")
	}
	storage := &db{
		cfg:       config,
		dbConnect: connect,
		log:       log,
	}

	if err := storage.dbConnect.Ping(); err != nil {
		log.Fatal("error connect to database")
	}

	return storage
}

func (db *db) logQuery(query Query) {
	db.log.Debug(query.String())
}
