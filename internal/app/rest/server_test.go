package rest

import (
	"testing"

	"github.com/stretchr/testify/suite"

	mock_service "github.com/training-of-new-employees/qon/mocks/service"
	mock_cache "github.com/training-of-new-employees/qon/mocks/store/cache"

	"github.com/golang/mock/gomock"
)

type handlerTestSuite struct {
	suite.Suite
	service         *mock_service.MockService
	userService     *mock_service.MockServiceUser
	positionService *mock_service.MockServicePosition
	lessonService   *mock_service.MockServiceLesson
	courseService   *mock_service.MockServiceCourse
	cache           *mock_cache.MockCache
	srv             *RestServer
}

// SetupSuite - запуск до начала выполнения набора тестов.
func (suite *handlerTestSuite) SetupSuite() {
}

// TearDownSuite - запуск после выполнения всех тестов.
func (suite *handlerTestSuite) TearDownSuite() {
}

// SetupTest - выполнение перед каждым тест-кейсом.
func (suite *handlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.userService = mock_service.NewMockServiceUser(ctrl)
	suite.positionService = mock_service.NewMockServicePosition(ctrl)
	suite.lessonService = mock_service.NewMockServiceLesson(ctrl)
	suite.courseService = mock_service.NewMockServiceCourse(ctrl)

	suite.service = mockService(
		ctrl,
		suite.userService,
		suite.positionService,
		suite.lessonService,
		suite.courseService,
	)
	suite.cache = mock_cache.NewMockCache(ctrl)

	suite.srv = New("secret", suite.service, suite.cache)
}

// TearDownTest - запуск после каждого тест-кейса.
func (suite *handlerTestSuite) TearDownTest() {
}

// TestStoreTestSuite - точка входа для тестирования хендлеров.
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(handlerTestSuite))
}

// mockService - инициализация моки-сервиса.
func mockService(ctrl *gomock.Controller, userService *mock_service.MockServiceUser, positionService *mock_service.MockServicePosition, lessonService *mock_service.MockServiceLesson, courseService *mock_service.MockServiceCourse) *mock_service.MockService {
	service := mock_service.NewMockService(ctrl)
	service.EXPECT().User().Return(userService).AnyTimes()
	service.EXPECT().Position().Return(positionService).AnyTimes()
	service.EXPECT().Lesson().Return(lessonService).AnyTimes()
	service.EXPECT().Course().Return(courseService).AnyTimes()
	return service
}
