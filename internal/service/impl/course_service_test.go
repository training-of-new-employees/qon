package impl

import (
	"context"
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/training-of-new-employees/qon/internal/model"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
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

func Test_courseService_GetCourses(t *testing.T) {
	type fields struct {
		coursedb *mock_store.MockRepositoryCourse
	}
	type args struct {
		ctx context.Context
		u   model.User
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
				model.User{
					ID:        1,
					CompanyID: 10,
					Name:      randomseq.RandomString(10),
					Email:     randomseq.RandomString(10),
					IsAdmin:   true,
				},
			},
			[]model.Course{},
			false,
		},
		{
			"Получение курсов пользователем",
			func(f *fields) {
				f.coursedb.EXPECT().UserCourses(gomock.Any(), 1).Return([]model.Course{}, nil)
			},
			args{
				nil,
				model.User{
					ID:        1,
					CompanyID: 10,
					Name:      randomseq.RandomString(10),
					Email:     randomseq.RandomString(10),
					IsAdmin:   false,
				},
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
			got, err := cs.GetCourses(tt.args.ctx, tt.args.u)
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
		want    model.Course
		wantErr bool
	}{
		{
			"Invalid course info",
			nil,
			args{
				nil,
				model.NewCourseSet(1, 10),
			},
			model.Course{},
			true,
		},
		{
			"Корректные данные курса",
			func(f *fields) {
				f.coursedb.EXPECT().CreateCourse(gomock.Any(), model.CourseSet{
					Name:        "Мой новый курс",
					Description: "Описание курса",
					CreatedBy:   10,
				}).Return(model.Course{
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
			model.Course{
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
		want    model.Course
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
			model.Course{},
			true,
		},
		{
			"Валидные данные",
			func(f *fields) {
				set := model.CourseSet{
					Name: "Новое имя",
				}
				ret := model.Course{
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
			model.Course{
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