package pg

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ store.RepositoryCourse = (*courseStorage)(nil)

type courseStorage struct {
	db *sqlx.DB
	transaction
}

func newCourseStorage(db *sqlx.DB, s *Store) *courseStorage {
	return &courseStorage{
		db:          db,
		transaction: transaction{db: db},
	}
}

func (c *courseStorage) UserCourses(ctx context.Context, userID int) ([]model.Course, error) {
	courses := make([]model.Course, 0, 10)
	qCourses := `SELECT c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
	FROM users u
	JOIN position_course pc ON u.position_id = pc.position_id
	JOIN courses c ON pc.course_id = c.id
	WHERE u.id = $1`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &courses, qCourses, &userID)
	})
	if err != nil {
		return nil, handleError(err)
	}
	if len(courses) == 0 {
		return nil, handleError(sql.ErrNoRows)
	}
	return courses, handleError(err)

}
func (c *courseStorage) CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	courses := make([]model.Course, 0, 10)
	qCourses := `SELECT c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
	FROM courses c
	JOIN users u ON u.id = c.created_by
	WHERE u.company_id = $1`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &courses, qCourses, &companyID)
	})
	if err != nil {
		return nil, handleError(err)
	}
	if len(courses) == 0 {
		return nil, handleError(sql.ErrNoRows)
	}
	return courses, handleError(err)
}

func (c *courseStorage) CreateCourse(ctx context.Context, course model.CourseSet) (*model.Course, error) {
	var res model.Course
	qCreate := `INSERT INTO courses (name, description, created_by)
	VALUES ($1, $2, $3)
	RETURNING id, created_by, active, archived, name, description, created_at, updated_at`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &res, qCreate, course.Name, course.Description, course.CreatedBy)
	})
	return &res, handleError(err)
}

func (c *courseStorage) EditCourse(ctx context.Context, course model.CourseSet, companyID int) (*model.Course, error) {
	var res model.Course
	qEdit := `UPDATE courses AS c
	SET name=COALESCE($1,c.name), description=COALESCE($2,c.description), archived=$3
	FROM users u  
	WHERE c.id=$4 AND u.company_id=$5 AND c.created_by = u.id
	RETURNING c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at`
	n := &course.Name
	d := &course.Description
	if course.Name == "" {
		n = nil
	}
	if course.Description == "" {
		d = nil
	}
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &res, qEdit, n, d, &course.IsArchived, &course.ID, &companyID)
	})
	return &res, handleError(err)

}
