package model

import (
	"testing"

	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func TestLesson_Validation(t *testing.T) {
	tests := []struct {
		name    string
		l       *Lesson
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "success",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30),
			},
			wantErr: nil,
		},
		{
			name: "empty lesson name",
			l: &Lesson{
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrLessonNameNotEmpty,
		},
		{
			name: "name contains smile",
			l: &Lesson{
				Name:    randomseq.RandomName(10) + "☺",
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "name contains *",
			l: &Lesson{
				Name:    randomseq.RandomName(10) + "*",
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "name contains special symbols",
			l: &Lesson{
				Name:    randomseq.RandomName(10) + "!№():,-?%'\";@",
				Content: randomseq.RandomString(30),
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			l: &Lesson{
				Name: randomseq.RandomName(10),
			},
			wantErr: errs.ErrTextContentNotEmpty,
		},
		{
			name: "too short content",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(minContentL - 1),
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "content contains special symbols",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(minContentL) + "!№():,-?%'\";@*#",
			},
			wantErr: nil,
		},
		{
			name: "too long content",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(maxContentL + 1),
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "valid url-picture",
			l: &Lesson{
				Name:       randomseq.RandomName(10),
				Content:    randomseq.RandomString(20),
				URLPicture: randomseq.RandomURLPicture(),
			},
			wantErr: nil,
		},
		{
			name: "invalid url-picture",
			l: &Lesson{
				Name:       randomseq.RandomName(10),
				Content:    randomseq.RandomString(20),
				URLPicture: randomseq.RandomString(maxURLPictureL + 1),
			},
			wantErr: errs.ErrInvalidURLPicture,
		},
		{
			name: "some url as url-picture",
			l: &Lesson{
				Name:       randomseq.RandomName(10),
				Content:    randomseq.RandomString(20),
				URLPicture: randomseq.RandomURL(),
			},
			wantErr: errs.ErrInvalidURLPicture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Validation(); err != tt.wantErr {
				t.Errorf("Lesson.Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLessonUpdate_Validation(t *testing.T) {
	tests := []struct {
		name    string
		l       *LessonUpdate
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "empty content",
			l: &LessonUpdate{
				Name: randomseq.RandomName(10),
			},
			wantErr: nil,
		},
		{
			name: "too long lesson name",
			l: &LessonUpdate{
				Name: randomseq.RandomName(maxLessonNameL + 1),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "content contains special symbols",
			l: &LessonUpdate{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30) + "!№():,-?%'\";@*#",
			},
			wantErr: nil,
		},
		{
			name: "content contains smile",
			l: &LessonUpdate{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30) + "☺",
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "valid url-picture",
			l: &LessonUpdate{
				Name:       randomseq.RandomName(10),
				URLPicture: randomseq.RandomURLPicture(),
			},
			wantErr: nil,
		},
		{
			name: "invalid url-picture",
			l: &LessonUpdate{
				Name:       randomseq.RandomName(10),
				URLPicture: randomseq.RandomString(minURLPictureL),
			},
			wantErr: errs.ErrInvalidURLPicture,
		},
		{
			name: "some url as url-picture",
			l: &LessonUpdate{
				Name:       randomseq.RandomName(10),
				URLPicture: randomseq.RandomURL(),
			},
			wantErr: errs.ErrInvalidURLPicture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Validation(); err != tt.wantErr {
				t.Errorf("LessonUpdate.Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
