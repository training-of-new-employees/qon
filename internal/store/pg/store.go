// Package pg - реализация интерфейса Store - работает с PostgeSQL.
package pg

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/training-of-new-employees/qon/internal/logger"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.Storages = (*Store)(nil)

// Store реализует интерфейс Store (для PostgreSQL).
type Store struct {
	conn          *sqlx.DB
	userStore     *uStorage
	positionStore *positionStorage
	companyStore  *companyStorage
	courseStorage *courseStorage
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
func (s *Store) Close() error {
	if err := s.conn.Close(); err != nil {
		logger.Log.Error("db close error", zap.Error(err))
		return err
	}
	logger.Log.Info("store closed successfully")

	return nil
}

// newPostgresDB устанавливает соединение с PostgreSQL.
func newPostgresDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("db open error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	return db, nil
}

// UserStorage - хранилище пользователей.
func (s *Store) UserStorage() store.RepositoryUser {
	if s.userStore != nil {
		return s.userStore
	}

	s.userStore = newUStorages(s.conn, s)
	return s.userStore
}

// PositionStorage - хранилище должностей.
func (s *Store) PositionStorage() store.RepositoryPosition {
	if s.positionStore != nil {
		return s.positionStore
	}

	s.positionStore = newPositionStorage(s.conn, s)

	return s.positionStore
}

// CompanyStorage - хранилище компаний.
func (s *Store) CompanyStorage() store.RepositoryCompany {
	if s.companyStore != nil {
		return s.companyStore
	}

	s.companyStore = newCompanyStorage(s.conn, s)

	return s.companyStore
}

// CourseStorage - хранилище курсов.
func (s *Store) CourseStorage() store.RepositoryCourse {
	if s.courseStorage != nil {
		return s.courseStorage
	}
	s.courseStorage = newCourseStorage(s.conn, s)
	return s.courseStorage
}
