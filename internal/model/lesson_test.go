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
			name: "успешная валидация",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30),
			},
			wantErr: nil,
		},
		{
			name: "пустое имя",
			l: &Lesson{
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrLessonNameNotEmpty,
		},
		{
			name: "слишком короткое имя",
			l: &Lesson{
				Name:    randomseq.RandomName(minNameL - 1),
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "смайл в имени",
			l: &Lesson{
				Name:    randomseq.RandomName(10) + "☺",
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "* в имени",
			l: &Lesson{
				Name:    randomseq.RandomName(10) + "*",
				Content: randomseq.RandomString(30),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "Имя со спец символами",
			l: &Lesson{
				Name:    randomseq.RandomName(10) + "!№():,-?%'\";@",
				Content: randomseq.RandomString(30),
			},
			wantErr: nil,
		},
		{
			name: "Пустое поле content",
			l: &Lesson{
				Name: randomseq.RandomName(10),
			},
			wantErr: errs.ErrTextContentNotEmpty,
		},
		{
			name: "Слишком короткое поле content",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(minContentL - 1),
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "# в поле content",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(minContentL) + "#",
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "Слишком длинное поле content",
			l: &Lesson{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(maxContentL + 1),
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "Слишком длинное поле url_picture",
			l: &Lesson{
				Name:       randomseq.RandomName(10),
				Content:    randomseq.RandomString(20),
				URLPicture: randomseq.RandomString(maxURLPictureL + 1),
			},
			wantErr: errs.ErrURLPictureLength,
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
			name: "Пустое поле content",
			l: &LessonUpdate{
				Name: randomseq.RandomName(10),
			},
			wantErr: nil,
		},
		{
			name: "слишком длинное имя",
			l: &LessonUpdate{
				Name: randomseq.RandomName(maxNameL + 1),
			},
			wantErr: errs.ErrInvalidLessonName,
		},
		{
			name: "* в поле content",
			l: &LessonUpdate{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30) + "*",
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "Content со спец символами",
			l: &LessonUpdate{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30) + "!№():,-?%'\";@",
			},
			wantErr: nil,
		},
		{
			name: "смайл в поле content",
			l: &LessonUpdate{
				Name:    randomseq.RandomName(10),
				Content: randomseq.RandomString(30) + "☺",
			},
			wantErr: errs.ErrInvalidTextContent,
		},
		{
			name: "Слишком короткое поле url_picture",
			l: &LessonUpdate{
				Name:       randomseq.RandomName(10),
				URLPicture: randomseq.RandomString(minURLPictureL - 1),
			},
			wantErr: errs.ErrURLPictureLength,
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
