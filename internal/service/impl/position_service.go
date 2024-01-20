package impl

import (
	"context"

	"github.com/training-of-new-employees/qon/internal/errs"
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

	position, err := p.db.PositionStorage().CreatePosition(ctx, val)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func (p *positionService) GetPosition(ctx context.Context, companyID int, positionID int) (*model.Position, error) {
	position, err := p.db.PositionStorage().GetPositionInCompany(ctx, companyID, positionID)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func (p *positionService) GetPositionCourses(ctx context.Context, companyID int, positionID int) ([]int, error) {
	_, err := p.db.PositionStorage().GetPositionInCompany(ctx, companyID, positionID)
	if err != nil {
		return nil, err
	}

	courseIDs, err := p.db.PositionStorage().GetCourseForPosition(ctx, positionID)
	if err != nil {
		return nil, err
	}

	return courseIDs, nil
}

func (p *positionService) GetPositions(ctx context.Context, id int) ([]*model.Position, error) {
	positions, err := p.db.PositionStorage().ListPositions(ctx, id)
	if err != nil {
		return nil, err
	}

	return positions, nil
}

func (p *positionService) UpdatePosition(ctx context.Context, id int, val model.PositionSet) (*model.Position, error) {
	position, err := p.db.PositionStorage().UpdatePosition(ctx, id, val)
	if err != nil {
		return nil, err
	}

	return position, nil
}

func (p *positionService) AssignCourse(ctx context.Context, positionID int, courseID int, userID int) error {
	user, err := p.db.UserStorage().GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	position, err := p.db.PositionStorage().GetPositionByID(ctx, positionID)
	if err != nil {
		return err
	}

	if user.CompanyID != position.CompanyID {
		return errs.ErrUnauthorized
	}

	if err := p.db.PositionStorage().AssignCourse(ctx, positionID, courseID); err != nil {
		return err
	}
	return nil
}

func (p *positionService) AssignCourses(ctx context.Context, positionID int, courseIDs []int, userID int) error {
	user, err := p.db.UserStorage().GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	position, err := p.db.PositionStorage().GetPositionByID(ctx, positionID)
	if err != nil {
		return err
	}

	if user.CompanyID != position.CompanyID {
		return errs.ErrUnauthorized
	}

	if err := p.db.PositionStorage().AssignCourses(ctx, positionID, courseIDs); err != nil {
		return err
	}
	return nil
}
