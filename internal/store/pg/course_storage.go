package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/training-of-new-employees/qon/internal/errs"
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
	qCourses := `
		SELECT
			c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
		FROM users u
		JOIN position_course pc ON u.position_id = pc.position_id
		JOIN courses c ON pc.course_id = c.id
		WHERE u.id = $1
	`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &courses, qCourses, &userID)
	})
	if err != nil {
		return nil, handleError(err)
	}

	return courses, nil
}

func (c *courseStorage) GetUserCourse(ctx context.Context, userID int, courseID int) (*model.Course, error) {
	courses := make([]model.Course, 0, 1)
	qCourses := `
		SELECT
			c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
		FROM users u
		JOIN position_course pc ON u.position_id = pc.position_id
		JOIN courses c ON pc.course_id = c.id
		WHERE u.id = $1 and c.id = $2
	`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &courses, qCourses, &userID, &courseID)
	})
	if err != nil {
		return nil, handleError(err)
	}
	if len(courses) == 0 {
		return nil, errs.ErrCourseNotFound
	}

	return &courses[0], nil
}

func (c *courseStorage) GetUserCoursesStatus(ctx context.Context, userID int, coursesIds []int) (map[int]string, error) {
	if len(coursesIds) == 0 {
		return map[int]string{}, nil
	}

	query := strings.Builder{}
	query.WriteString(`INSERT INTO course_assign (user_id, course_id) VALUES `)

	var params []interface{}

	for i, courseID := range coursesIds {
		position := i * 2
		query.WriteString(fmt.Sprintf("($%d,$%d)", position+1, position+2))
		params = append(params, userID, courseID)

		if i+1 < len(coursesIds) {
			query.WriteString(",")
		}
	}

	query.WriteString(" ON CONFLICT (user_id, course_id) DO UPDATE SET user_id = EXCLUDED.user_id RETURNING course_id, status")
	queryStr := query.String()
	statuses := make(map[int]string)

	err := c.tx(func(tx *sqlx.Tx) error {
		for _, courseID := range coursesIds {
			if err := c.transaction.syncUserCourseProgressTx(
				ctx,
				tx,
				userID,
				courseID,
			); err != nil {
				return err
			}
		}

		rows, err := tx.QueryContext(ctx, queryStr, params...)
		if err != nil {
			return err
		}

		for rows.Next() {
			var courseID int
			var status string

			if err := rows.Scan(&courseID, &status); err != nil {
				return err
			}

			statuses[courseID] = status
		}

		return nil
	})

	if err != nil {
		return nil, handleError(err)
	}

	return statuses, nil
}

func (c *courseStorage) CompanyCourses(ctx context.Context, companyID int) ([]model.Course, error) {
	courses := make([]model.Course, 0, 10)
	qCourses := `
		SELECT
			c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
		FROM courses c
		JOIN users u ON u.id = c.created_by
		WHERE u.company_id = $1
	`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &courses, qCourses, &companyID)
	})
	if err != nil {
		return nil, handleError(err)
	}

	return courses, nil
}

func (c *courseStorage) CompanyCourse(ctx context.Context, courseID, companyID int) (*model.Course, error) {
	qCourse := `SELECT c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at 
	FROM courses c
	JOIN users u ON u.id = c.created_by
	WHERE u.company_id = $1 AND c.id=$2`
	var course model.Course
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.SelectContext(ctx, &course, qCourse, &companyID, &courseID)
	})
	if err != nil {
		return nil, handleError(err)
	}

	return &course, nil
}

func (c *courseStorage) CreateCourse(ctx context.Context, course model.CourseSet) (*model.Course, error) {
	var res model.Course
	qCreate := `
		INSERT INTO courses (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_by, active, archived, name, description, created_at, updated_at
	`
	err := c.tx(func(tx *sqlx.Tx) error {
		return tx.GetContext(ctx, &res, qCreate, course.Name, course.Description, course.CreatedBy)
	})
	return &res, handleError(err)
}

func (c *courseStorage) EditCourse(ctx context.Context, course model.CourseSet, companyID int) (*model.Course, error) {
	var res model.Course
	qEdit := `
		UPDATE courses AS c
		SET name=COALESCE($1,c.name), description=COALESCE($2,c.description), archived=$3
		FROM users u  
		WHERE c.id=$4 AND u.company_id=$5 AND c.created_by = u.id
		RETURNING c.id, c.created_by, c.active, c.archived, c.name, c.description, c.created_at, c.updated_at
	`
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
