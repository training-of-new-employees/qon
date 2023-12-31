package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

// RepositoryPosition - интерфейс репозитория должности.
type RepositoryPosition interface {
	// CreatePositionDB - создание должности
	CreatePositionDB(ctx context.Context, position model.PositionSet) (*model.Position, error)
	// GetPositionDB - получение данных должности, привязанной к компании
	// TODO: возможно этот метод стоит убрать, вместо можно использовать GetPositionByID
	GetPositionDB(ctx context.Context, companyID int, positionID int) (*model.Position, error)
	// GetPositionsDB - получение должностей компании
	GetPositionsDB(ctx context.Context, companyID int) ([]*model.Position, error)
	// GetPositionByID - получение данных должности по ID
	GetPositionByID(ctx context.Context, positionID int) (*model.Position, error)
	// UpdatePositionDB - обновление данных должности
	UpdatePositionDB(ctx context.Context, id int, position model.PositionSet) (*model.Position, error)
	// AssignCourseDB - назначение курса на должность
	AssignCourseDB(ctx context.Context, positionID int, courseID int, user_id int) error
}
