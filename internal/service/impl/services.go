package impl

import (
	"time"

	"github.com/training-of-new-employees/qon/internal/pkg/doar"
	"github.com/training-of-new-employees/qon/internal/pkg/jwttoken"
	"github.com/training-of-new-employees/qon/internal/service"
	"github.com/training-of-new-employees/qon/internal/store"
	"github.com/training-of-new-employees/qon/internal/store/cache"
)

// Services - структура, которая содержит в себе все сервисы.
type Services struct {
	db              store.Storages
	cache           cache.Cache
	secretKey       string
	aTokenTime      time.Duration
	rTokenTime      time.Duration
	userService     *uService
	positionService *positionService
	courseService   *courseService
	lessonService   *lessonService
	sender          doar.EmailSender
	host            string
}

func NewServices(db store.Storages, cache cache.Cache, secretKey string, aTokTimeDur time.Duration, rTokTimeDur time.Duration, sender doar.EmailSender, host string) *Services {
	return &Services{
		db:         db,
		cache:      cache,
		secretKey:  secretKey,
		aTokenTime: aTokTimeDur,
		rTokenTime: rTokTimeDur,
		sender:     sender,
		host:       host,
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
		s.sender,
		s.host,
	)

	return s.userService
}

func (s *Services) Position() service.ServicePosition {

	if s.positionService != nil {
		return s.positionService
	}

	s.positionService = newPositionService(s.db)

	return s.positionService
}

func (s *Services) Course() service.ServiceCourse {
	if s.courseService != nil {
		return s.courseService
	}
	s.courseService = newCourseService(s.db)
	return s.courseService
}

func (s *Services) Lesson() service.ServiceLesson {
	if s.lessonService != nil {
		return s.lessonService
	}

	s.lessonService = newLessonService(s.db)

	return s.lessonService
}
