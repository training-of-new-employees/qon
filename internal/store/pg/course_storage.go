package pg

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryCourse = (*courseStorage)(nil)

type courseStorage struct {
	db    *sqlx.DB
	store *Store
}

func newCourseStorage(db *sqlx.DB, s *Store) *courseStorage {
	return &courseStorage{db: db, store: s}
}

func (c *courseStorage) PositionCourses(ctx context.Context, userID int) ([]model.Course, error) {
	return nil, nil

}
func (c *courseStorage) CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	return nil, nil
}
