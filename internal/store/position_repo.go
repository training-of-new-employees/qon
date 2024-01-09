package store

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/model"
)

// RepositoryPosition - интерфейс репозитория должности.
type RepositoryPosition interface {
	// TODO: Не совсем понятно, зачем в методах используется окончание 'DB', возможно, его следует потом убрать.
	CreatePositionDB(ctx context.Context, position model.PositionSet) (*model.Position, error)

	// GetPositionDB - получение данных должности, привязанной к компании
	// TODO: возможно этот метод стоит убрать, вместо можно использовать GetPositionByID
	GetPositionDB(ctx context.Context, companyID int, positionID int) (*model.Position, error)

	// TODO: Названия методов GetPositionDB и GetPositionsDB создают путаницу, т.к. раличаются только наличием одной буквы 's'.
	// Возможно, следует изменить название метода на ListPositionDB (для примера).
	GetPositionsDB(ctx context.Context, companyID int) ([]*model.Position, error)

	GetPositionByID(ctx context.Context, positionID int) (*model.Position, error)
	UpdatePositionDB(ctx context.Context, positionID int, position model.PositionSet) (*model.Position, error)
	AssignCourseDB(ctx context.Context, positionID int, courseID int) error
}
