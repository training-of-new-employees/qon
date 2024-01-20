package impl

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	mock_doar "github.com/training-of-new-employees/qon/mocks/pkg/doar"
	mock_jwttoken "github.com/training-of-new-employees/qon/mocks/pkg/jwttoken"
	mock_store "github.com/training-of-new-employees/qon/mocks/store"
	mock_cache "github.com/training-of-new-employees/qon/mocks/store/cache"

	"go.uber.org/mock/gomock"
)

type serviceTestSuite struct {
	suite.Suite

	userService     *uService
	lessonService   *lessonService
	courseService   *courseService
	positionService *positionService

	jwtgenerator *mock_jwttoken.MockJWTGenerator
	jwtvalidator *mock_jwttoken.MockJWTValidator
	emailSender  *mock_doar.MockEmailSender
	cache        *mock_cache.MockCache

	db              *mock_store.MockStorages
	userStorage     *mock_store.MockRepositoryUser
	positionStorage *mock_store.MockRepositoryPosition
	courseStorage   *mock_store.MockRepositoryCourse
	companyStorage  *mock_store.MockRepositoryCompany
	lessonStorage   *mock_store.MockRepositoryLesson

	secret    string
	host      string
	aTokenDur time.Duration
	rTokenDur time.Duration

	service *Services
}

// SetupSuite - запуск до начала выполнения набора тестов.
func (suite *serviceTestSuite) SetupSuite() {
}

// TearDownSuite - запуск после выполнения всех тестов.
func (suite *serviceTestSuite) TearDownSuite() {
}

// SetupTest - выполнение перед каждым тест-кейсом.
func (suite *serviceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.userStorage = mock_store.NewMockRepositoryUser(ctrl)
	suite.positionStorage = mock_store.NewMockRepositoryPosition(ctrl)
	suite.lessonStorage = mock_store.NewMockRepositoryLesson(ctrl)
	suite.courseStorage = mock_store.NewMockRepositoryCourse(ctrl)
	suite.companyStorage = mock_store.NewMockRepositoryCompany(ctrl)

	suite.db = mockDB(
		ctrl,
		suite.userStorage,
		suite.positionStorage,
		suite.lessonStorage,
		suite.companyStorage,
		suite.courseStorage,
	)
	suite.cache = mock_cache.NewMockCache(ctrl)
	suite.jwtgenerator = mock_jwttoken.NewMockJWTGenerator(ctrl)
	suite.jwtvalidator = mock_jwttoken.NewMockJWTValidator(ctrl)
	suite.emailSender = mock_doar.NewMockEmailSender(ctrl)
	suite.secret = "secret"
	suite.host = "test host"
	suite.rTokenDur = time.Hour * 1
	suite.aTokenDur = time.Hour * 1

	suite.userService = newUserService(
		suite.db,
		suite.secret,
		suite.aTokenDur,
		suite.rTokenDur,
		suite.cache,
		suite.jwtgenerator,
		suite.jwtvalidator,
		suite.emailSender,
		suite.host,
	)
	suite.positionService = newPositionService(suite.db)
	suite.lessonService = newLessonService(suite.db)
	suite.courseService = newCourseService(suite.db)

	suite.service = NewServices(
		suite.db,
		suite.cache,
		suite.secret,
		suite.aTokenDur,
		suite.rTokenDur,
		suite.emailSender,
		suite.host,
	)
}

// TearDownTest - запуск после каждого тест-кейса.
func (suite *serviceTestSuite) TearDownTest() {
}

// TestServiceTestSuite - точка входа для тестирования сервисов.
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(serviceTestSuite))
}

// mockDB - инициализация моки-бд.
func mockDB(
	ctrl *gomock.Controller,
	userRepository *mock_store.MockRepositoryUser,
	positionRepository *mock_store.MockRepositoryPosition,
	lessonRepository *mock_store.MockRepositoryLesson,
	companyRepository *mock_store.MockRepositoryCompany,
	courseRepository *mock_store.MockRepositoryCourse,
) *mock_store.MockStorages {
	db := mock_store.NewMockStorages(ctrl)
	db.EXPECT().PositionStorage().Return(positionRepository).AnyTimes()
	db.EXPECT().UserStorage().Return(userRepository).AnyTimes()
	db.EXPECT().LessonStorage().Return(lessonRepository).AnyTimes()
	db.EXPECT().CourseStorage().Return(courseRepository).AnyTimes()
	db.EXPECT().CompanyStorage().Return(companyRepository).AnyTimes()
	return db
}
