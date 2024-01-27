package impl

import (
	"context"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/model"
	mock_store "github.com/training-of-new-employees/qon/mocks/store"
)

func Test_newCourseService(t *testing.T) {
	tests := []struct {
		name string
		want *courseService
	}{{
		"Create empty storage",
		&courseService{},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mock_store.NewMockStorages(ctrl)
			if got := newCourseService(s); got == nil {
				t.Errorf("newCourseService() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_courseService_GetAdminCourses(t *testing.T) {
	type fields struct {
		coursedb *mock_store.MockRepositoryCourse
	}
	type args struct {
		ctx       context.Context
		companyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    []model.Course
		wantErr bool
	}{
		{
			"Получение курсов админом",
			func(f *fields) {
				f.coursedb.EXPECT().CompanyCourses(gomock.Any(), 10).Return([]model.Course{}, nil)
			},
			args{
				nil,
				10,
			},
			[]model.Course{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cdb := mock_store.NewMockRepositoryCourse(ctrl)
			f := &fields{coursedb: cdb}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			storage := mockCourseStorage(ctrl, f.coursedb)
			cs := &courseService{
				db: storage,
			}
			got, err := cs.GetCompanyCourses(tt.args.ctx, tt.args.companyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("courseService.GetCourses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("courseService.GetCourses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_courseService_CreateCourse(t *testing.T) {
	type fields struct {
		coursedb *mock_store.MockRepositoryCourse
	}
	type args struct {
		ctx context.Context
		c   model.CourseSet
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.Course
		wantErr bool
	}{
		{
			"Корректные данные курса",
			func(f *fields) {
				f.coursedb.EXPECT().CreateCourse(gomock.Any(), model.CourseSet{
					Name:        "Мой новый курс",
					Description: "Описание курса",
					CreatedBy:   10,
				}).Return(&model.Course{
					ID:          1,
					CreatedBy:   10,
					IsActive:    true,
					Name:        "Мой новый курс",
					Description: "Описание курса",
					IsArchived:  false,
				}, nil)

			},
			args{
				nil,
				model.CourseSet{
					Name:        "Мой новый курс",
					Description: "Описание курса",
					CreatedBy:   10,
				},
			},
			&model.Course{
				ID:          1,
				CreatedBy:   10,
				IsActive:    true,
				Name:        "Мой новый курс",
				Description: "Описание курса",
				IsArchived:  false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cdb := mock_store.NewMockRepositoryCourse(ctrl)
			f := &fields{
				cdb,
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			cs := &courseService{
				db: mockCourseStorage(ctrl, f.coursedb),
			}
			got, err := cs.CreateCourse(tt.args.ctx, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("courseService.CreateCourse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("courseService.CreateCourse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_courseService_EditCourse(t *testing.T) {
	type fields struct {
		courseDB *mock_store.MockRepositoryCourse
	}
	type args struct {
		ctx       context.Context
		c         model.CourseSet
		companyID int
	}
	tests := []struct {
		name    string
		prepare func(*fields)
		args    args
		want    *model.Course
		wantErr bool
	}{
		{
			"Некорректные данные",
			nil,
			args{
				nil,
				model.CourseSet{
					Name:        "n",
					Description: "",
					IsArchived:  true,
				},
				1,
			},
			nil,
			true,
		},
		{
			"Валидные данные",
			func(f *fields) {
				set := model.CourseSet{
					Name: "Новое имя",
				}
				ret := &model.Course{
					ID:         1,
					Name:       set.Name,
					IsArchived: false,
				}
				f.courseDB.EXPECT().EditCourse(gomock.Any(), set, 10).Return(ret, nil)
			},
			args{
				nil,
				model.CourseSet{
					Name: "Новое имя",
				},
				10,
			},
			&model.Course{
				ID:         1,
				Name:       "Новое имя",
				IsArchived: false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			f := &fields{
				mock_store.NewMockRepositoryCourse(ctrl),
			}
			if tt.prepare != nil {
				tt.prepare(f)
			}
			cs := &courseService{
				db: mockCourseStorage(ctrl, f.courseDB),
			}
			got, err := cs.EditCourse(tt.args.ctx, tt.args.c, tt.args.companyID)
			if (err != nil) != tt.wantErr {
				t.Errorf("courseService.EditCourse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("courseService.EditCourse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockCourseStorage(ctrl *gomock.Controller, cStore *mock_store.MockRepositoryCourse) *mock_store.MockStorages {
	storages := mock_store.NewMockStorages(ctrl)
	storages.EXPECT().CourseStorage().Return(cStore).AnyTimes()
	return storages
}

func (suite *serviceTestSuite) TestGetUserCourseLessons() {
	testCases := []struct {
		name    string
		err     error
		prepare func() (int, int)
	}{
		{
			name: "success",
			err:  nil,
			prepare: func() (int, int) {
				courseID := 1
				userID := 1

				suite.courseStorage.
					EXPECT().
					GetUserCourse(gomock.Any(), userID, courseID).
					Return(&model.Course{ID: courseID}, nil)

				suite.lessonStorage.
					EXPECT().
					GetLessonsList(gomock.Any(), courseID).
					Return([]model.Lesson{{ID: 1}}, nil)

				suite.lessonStorage.
					EXPECT().
					GetUserLessonsStatus(gomock.Any(), userID, courseID, []int{1}).
					Return(map[int]string{1: "not-started"}, nil)

				return courseID, userID
			},
		},
		{
			name: "course not found",
			err:  errs.ErrCourseNotFound,
			prepare: func() (int, int) {
				courseID := 1
				userID := 1

				suite.courseStorage.
					EXPECT().
					GetUserCourse(gomock.Any(), userID, courseID).
					Return(nil, errs.ErrCourseNotFound)

				return courseID, userID
			},
		},
		{
			name: "error getting lesson",
			err:  errs.ErrInternal,
			prepare: func() (int, int) {
				courseID := 1
				userID := 1

				suite.courseStorage.
					EXPECT().
					GetUserCourse(gomock.Any(), userID, courseID).
					Return(&model.Course{ID: courseID}, nil)

				suite.lessonStorage.
					EXPECT().
					GetLessonsList(gomock.Any(), courseID).
					Return([]model.Lesson{{ID: 1}}, errs.ErrInternal)

				return courseID, userID
			},
		},
		{
			name: "error getting statuses",
			err:  errs.ErrInternal,
			prepare: func() (int, int) {
				courseID := 1
				userID := 1

				suite.courseStorage.
					EXPECT().
					GetUserCourse(gomock.Any(), userID, courseID).
					Return(&model.Course{ID: courseID}, nil)

				suite.lessonStorage.
					EXPECT().
					GetLessonsList(gomock.Any(), courseID).
					Return([]model.Lesson{{ID: 1}}, nil)

				suite.lessonStorage.
					EXPECT().
					GetUserLessonsStatus(gomock.Any(), userID, courseID, []int{1}).
					Return(nil, errs.ErrInternal)

				return courseID, userID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			courseID, userID := tc.prepare()
			_, err := suite.courseService.GetUserCourseLessons(context.TODO(), userID, courseID)
			suite.Equal(tc.err, err)
		})
	}
}

func (suite *serviceTestSuite) TestGetUserCourses() {
	testCases := []struct {
		name    string
		err     error
		prepare func() int
	}{
		{
			name: "success",
			err:  nil,
			prepare: func() int {
				userID := 1

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return([]model.Course{{ID: 1}}, nil)

				suite.courseStorage.
					EXPECT().
					GetUserCoursesStatus(gomock.Any(), userID, []int{1}).
					Return(map[int]string{1: "not-started"}, nil)

				return userID
			},
		},
		{
			name: "success (no courses)",
			err:  nil,
			prepare: func() int {
				userID := 1

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return([]model.Course{}, nil)

				return userID
			},
		},
		{
			name: "course internal error",
			err:  errs.ErrInternal,
			prepare: func() int {
				userID := 1

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return(nil, errs.ErrInternal)

				return userID
			},
		},
		{
			name: "course status internal error",
			err:  errs.ErrInternal,
			prepare: func() int {
				userID := 1

				suite.courseStorage.
					EXPECT().
					UserCourses(gomock.Any(), userID).
					Return([]model.Course{{ID: 1}}, nil)

				suite.courseStorage.
					EXPECT().
					GetUserCoursesStatus(gomock.Any(), userID, []int{1}).
					Return(nil, errs.ErrInternal)

				return userID
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			userID := tc.prepare()
			_, err := suite.courseService.GetUserCourses(context.TODO(), userID)
			suite.Equal(tc.err, err)
		})
	}
}
