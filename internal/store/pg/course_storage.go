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
	qPos := `SELECT position_id FROM users WHERE id = $1`
	qCoursesID := `SELECT course_id FROM position_course WHERE position_id = $1`
	qCourses, err := c.db.PreparexContext(ctx, `SELECT id, created_by, active, archived, name, description, created_at, updated_at FROM courses WHERE id = $1`)
	if err != nil {
		return nil, handleError(err)
	}

	err = tx(c.db, func(tx *sqlx.Tx) error {
		var posID int
		err := tx.GetContext(ctx, &posID, qPos, userID)
		if err != nil {
			return err
		}
		var coursesID []int
		err = tx.SelectContext(ctx, &coursesID, qCoursesID, posID)
		if err != nil {
			return err
		}
		courses = make([]model.Course, len(coursesID))
		qCourses := tx.Stmtx(qCourses)
		for i, id := range coursesID {
			err := qCourses.Get(&courses[i], id)
			if err != nil {
				return err
			}
		}
		return nil

	})
	return courses, err

}
func (c *courseStorage) CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	return nil, nil
}
