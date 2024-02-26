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
			"max length name",
			fields{
				Name: randomseq.RandomString(maxCourseNameL),
			},
			nil,
		},
		{
			"too long name",
			fields{
				Name: randomseq.RandomString(maxCourseNameL + 1),
			},
			errs.ErrInvalidCourseName,
		},
		{
			"name with *",
			fields{
				Name: randomseq.RandomString(minCourseNameL) + "*",
			},
			errs.ErrInvalidCourseName,
		},
		{
			"name with #",
			fields{
				Name: randomseq.RandomString(minCourseNameL) + "#",
			},
			errs.ErrInvalidCourseName,
		},
		{
			"name with special symbols",
			fields{
				Name: randomseq.RandomString(minCourseNameL) + "!№():,-?%'\";@",
			},
			nil,
		},
		{
			"name with punctuation (,)",
			fields{
				Name: strings.Join([]string{randomseq.RandomString(minCourseNameL), randomseq.RandomString(minCourseNameL)}, ","),
			},
			nil,
		},
		{
			"name with punctuation (;)",
			fields{
				Name: strings.Join([]string{randomseq.RandomString(minCourseNameL), randomseq.RandomString(minCourseNameL)}, ";"),
			},
			nil,
		},
		{
			"name contains space",
			fields{
				Name: strings.Join([]string{randomseq.RandomString(minCourseNameL), randomseq.RandomString(minCourseNameL)}, " "),
			},
			nil,
		},
		{
			"name contains smile",
			fields{
				Name: randomseq.RandomString(minCourseNameL) + "☺",
			},
			errs.ErrInvalidCourseName,
		},
		{
			"empty course name",
			fields{},
			errs.ErrCourseNameIsEmpty,
		},
		{
			"description contains *#",
			fields{
				Name:        randomseq.RandomString(minCourseNameL),
				Description: randomseq.RandomString(minCourseDescL) + "*#",
			},
			nil,
		},
		{
			"too short description",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(minCourseDescL - 1),
			},
			errs.ErrInvalidCourseDescription,
		},
		{
			"min length description",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(minCourseDescL),
			},
			nil,
		},
		{
			"max length description",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(maxCourseDescL),
			},
			nil,
		},
		{
			"max description with space",
			fields{
				Name:        "validname",
				Description: strings.Join([]string{randomseq.RandomString(maxCourseDescL / 2), randomseq.RandomString(maxCourseDescL/2 + 1)}, " "),
			},
			errs.ErrInvalidCourseDescription,
		},
		{
			"too long description",
			fields{
				Name:        "validname",
				Description: randomseq.RandomString(maxCourseDescL + 1),
			},
			errs.ErrInvalidCourseDescription,
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
			name: "проверка вызова валидации курса",
			fields: fields{
				Name: randomseq.RandomString(4),
			},
			wantErr: false,
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
			name: "creating course",
			args: args{0, 10},
			want: CourseSet{
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
