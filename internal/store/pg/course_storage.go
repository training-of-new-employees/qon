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
	var courses []model.Course
	qCourses := `SELECT c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
	FROM users u
	JOIN position_course pc ON u.position_id = pc.position_id
	JOIN courses c ON pc.course_id = c.id
	WHERE u.id = $1`
	err := tx(c.db, func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, courses, qCourses, userID)
	})
	return courses, err

}
func (c *courseStorage) CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	var courses []model.Course
	qCourses := `SELECT c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
	FROM positions p
	JOIN position_course pc ON p.id = pc.position_id
	JOIN courses c ON pc.course_id = c.id
	WHERE p.company_id = $1`
	err := tx(c.db, func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, courses, qCourses, companyID)
	})
	return courses, err
}

func (c *courseStorage) CreateCourse(ctx context.Context, course model.CourseSet, creatorID int) (*model.Course, error) {
	res := &model.Course{}
	qCreate := `INSERT INTO courses (name, description, created_by) VALUES ($1, $2, $3) RETURNING id, created_by, active, archived, name, description, created_at, updated_at`
	err := tx(c.db, func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, res, qCreate, course.Name, course.Description, creatorID)
	})
	return res, err
}
