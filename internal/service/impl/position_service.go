package impl

import (
	"context"
	"fmt"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
)

var _ service.ServicePosition = (*positionService)(nil)

type positionService struct {
	db store.Storages
}

func newPositionService(db store.Storages) *positionService {
	return &positionService{db: db}
}

func (p *positionService) CreatePosition(ctx context.Context, val model.PositionSet) (*model.Position, error) {

	position, err := p.db.PositionStorage().CreatePositionDB(ctx, val)
	if err != nil {
		return nil, fmt.Errorf("failed CreatePositionDB: %w", err)
	}

	return position, nil
}

func (p *positionService) GetPosition(ctx context.Context, companyID int, positionID int) (*model.Position, error) {
	position, err := p.db.PositionStorage().GetPositionDB(ctx, companyID, positionID)
	if err != nil {
		return nil, fmt.Errorf("failed GetPositionDB: %w", err)
	}

	return position, nil
}

func (p *positionService) GetPositions(ctx context.Context, id int) ([]*model.Position, error) {
	positions, err := p.db.PositionStorage().GetPositionsDB(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed GetPositionsDB: %w", err)
	}

	return positions, nil
}

func (p *positionService) UpdatePosition(ctx context.Context, id int, val model.PositionSet) (*model.Position, error) {
	position, err := p.db.PositionStorage().UpdatePositionDB(ctx, id, val)
	if err != nil {
		return nil, fmt.Errorf("failed UpdatePositionDB: %w", err)
	}

	return position, nil
}

func (p *positionService) DeletePosition(ctx context.Context, id int, companyID int) error {
	err := p.db.PositionStorage().DeletePositionDB(ctx, id, companyID)
	if err != nil {
		return fmt.Errorf("failed DeletePositionDB: %w", err)
	}

	return nil
}

func (p *positionService) AssignCourse(ctx context.Context, positionID int,
	courseID int, user_id int) error {
	user, err := p.db.UserStorage().GetUserByID(ctx, user_id)
	if err != nil {
		return err
	}
	position, err := p.db.PositionStorage().GetPositionByID(ctx, positionID)
	if err != nil {
		return err
	}

	if user.CompanyID != position.CompanyID {
		return model.ErrNoAuthorized
	}

	if err := p.db.PositionStorage().AssignCourseDB(ctx,
		positionID, courseID, user_id); err != nil {
		return err
	}
	return nil
}
