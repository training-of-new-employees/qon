package impl

import (
	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
	"time"
)

// Services - структура, которая содержит в себе все сервисы.
type Services struct {
	db          store.Storages
	cache       cache.Cache
	secretKey   string
	aTokenTime  time.Duration
	rTokenTime  time.Duration
	userService *uService
	sender      doar.EmailSender
}

func NewServices(db store.Storages, cache cache.Cache, secretKey string, aTokTimeDur time.Duration, rTokTimeDur time.Duration) *Services {
	return &Services{
		db:         db,
		cache:      cache,
		secretKey:  secretKey,
		aTokenTime: aTokTimeDur,
		rTokenTime: rTokTimeDur,
	}
}

func (s *Services) User() service.ServiceUser {

	if s.userService != nil {
		return s.userService
	}

	s.userService = newUserService(
		s.db,
		s.secretKey,
		s.aTokenTime,
		s.rTokenTime,
		s.cache,
		jwttoken.NewTokenGenerator(s.secretKey),
		jwttoken.NewTokenValidator(s.secretKey),
	)

	return s.userService
}
