// Package pg - реализация интерфейса Store - работает с PostgeSQL.
package pg

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/training-of-new-employees/qon/internal/logger"
	"go.uber.org/zap"
)

// Store реализует интерфейс Store (для PostgreSQL).
type Store struct {
	conn *sql.DB
}

// NewStore - конструктор для Store.
func NewStore(dsn string) (*Store, error) {
	// create connection to db
	db, err := newPostgresDB(dsn)
	if err != nil {
		return nil, err
	}
	logger.Log.Info("connection to db established")

	if err := MigrationsUp(db); err != nil {
		return nil, err
	}
	logger.Log.Info("db migrated")

	s := &Store{
		conn: db,
	}
	logger.Log.Info("store successfully created")

	return s, nil
}

// Close - деструктор для store.
func (s *Store) Close() {
	if err := s.conn.Close(); err != nil {
		logger.Log.Error("db close error", zap.Error(err))
		return
	}
	logger.Log.Info("store closed successfully")
}

// newPostgresDB устанавливает соединение с PostgreSQL.
func newPostgresDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	return db, nil
}
