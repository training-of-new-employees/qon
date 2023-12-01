// Package pg - реализация интерфейса Store - работает с PostgeSQL.
package pg

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
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
	log.Println("connection to db established")

	s := &Store{
		conn: db,
	}
	log.Println("store successfully created")

	return s, nil
}

// Close - деструктор для store.
func (s *Store) Close() {
	if err := s.conn.Close(); err != nil {
		log.Println("db close error: ", err)
		return
	}
	log.Println("store has been closed successfully")
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
