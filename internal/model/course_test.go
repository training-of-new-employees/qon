package model

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func TestCourse_Validation(t *testing.T) {
	type fields struct {
		ID          int
		CreatedBy   int
		IsActive    bool
		IsArchived  bool
		Name        string
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			"Короткое название",
			fields{
				Name: randomseq.RandomString(minNameL - 1),
			},
			errs.ErrCourseNameInvalid,
		},
		{
			"Максимально допустимая длина названия",
			fields{
				Name: randomseq.RandomString(maxNameL),
			},
			nil,
		},
		{
			"Слишком длинное название",
			fields{
				Name: randomseq.RandomString(maxNameL + 1),
			},
			errs.ErrCourseNameInvalid,
		},
		{
			"Название с символом *",
			fields{
				Name: randomseq.RandomString(minNameL) + "*",
			},
			errs.ErrCourseNameInvalid,
		},
		{
			"Название с символом #",
			fields{
				Name: randomseq.RandomString(minNameL) + "#",
			},
			errs.ErrCourseNameInvalid,
		},
		{
			"Название со спец символами",
			fields{
				Name: randomseq.RandomString(minNameL) + "!№():,-?%'\";@",
			},
			nil,
		},
		{
			"Название со знаками препинания",
			fields{
				Name: strings.Join([]string{randomseq.RandomString(minNameL), randomseq.RandomString(minNameL)}, ","),
			},
			nil,
		},
		{
			"Название со знаками препинания",
			fields{
				Name: strings.Join([]string{randomseq.RandomString(minNameL), randomseq.RandomString(minNameL)}, ";"),
			},
			nil,
		},
		{
			"Название с пробелом",
			fields{
				Name: strings.Join([]string{randomseq.RandomString(minNameL), randomseq.RandomString(minNameL)}, " "),
			},
			nil,
		},
		{
			"Название со смайлом",
			fields{
				Name: randomseq.RandomString(minNameL) + "☺",
			},
			errs.ErrCourseNameInvalid,
		},
		{
			"Пустое название",
			fields{},
			errs.ErrCourseNameIsEmpty,
		},
		{
			"Описание с символом #",
			fields{
				Name: randomseq.RandomString(minDescL) + "#",
			},
			errs.ErrCourseNameInvalid,
		},
		{
			"Слишком короткое описание",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(minDescL - 1),
			},
			errs.ErrCourseDescriptionInvalid,
		},
		{
			"Максимально короткое описание",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(minDescL),
			},
			nil,
		},
		{
			"Максимально длинное описание",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(maxDescL),
			},
			nil,
		},
		{
			"Максимально длинное описание с пробелом",
			fields{
				Name:        "validname",
				Description: strings.Join([]string{randomseq.RandomString(maxDescL / 2), randomseq.RandomString(maxDescL/2 + 1)}, " "),
			},
			errs.ErrCourseDescriptionInvalid,
		},
		{
			"Слишком длинное описание",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(maxDescL + 1),
			},
			errs.ErrCourseDescriptionInvalid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Course{
				ID:          tt.fields.ID,
				CreatedBy:   tt.fields.CreatedBy,
				IsActive:    tt.fields.IsActive,
				IsArchived:  tt.fields.IsArchived,
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				CreatedAt:   tt.fields.CreatedAt,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if err := c.Validation(); err != tt.wantErr {
				t.Errorf("Course.Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCourseSet_Validation(t *testing.T) {
	type fields struct {
		ID          int
		CreatedBy   int
		Name        string
		Description string
		IsArchived  bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Проверка вызова валидации курса",
			fields{
				Name: randomseq.RandomString(4),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &CourseSet{
				ID:          tt.fields.ID,
				CreatedBy:   tt.fields.CreatedBy,
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				IsArchived:  tt.fields.IsArchived,
			}
			if err := cs.Validation(); (err != nil) != tt.wantErr {
				t.Errorf("CourseSet.Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewCourseSet(t *testing.T) {
	type args struct {
		id      int
		creator int
	}
	tests := []struct {
		name string
		args args
		want CourseSet
	}{
		{
			"Создание курса",
			args{0, 10},
			CourseSet{
				ID:        0,
				CreatedBy: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCourseSet(tt.args.id, tt.args.creator); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCourseSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
