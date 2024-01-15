package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

// RepositoryPosition - интерфейс репозитория должности.
type RepositoryPosition interface {
	CreatePosition(ctx context.Context, position model.PositionSet) (*model.Position, error)

	GetPositionByID(ctx context.Context, positionID int) (*model.Position, error)

	// GetPositionInComp - получение данных должности, привязанной к компании
	GetPositionInComp(ctx context.Context, companyID int, positionID int) (*model.Position, error)

	ListPositions(ctx context.Context, companyID int) ([]*model.Position, error)

	UpdatePosition(ctx context.Context, positionID int, position model.PositionSet) (*model.Position, error)

	AssignCourse(ctx context.Context, positionID int, courseID int) error
}
