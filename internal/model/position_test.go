package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/training-of-new-employees/qon/internal/errs"
	"github.com/training-of-new-employees/qon/internal/pkg/randomseq"
)

func TestPositionEdit_Validation(t *testing.T) {
	type fields struct {
		CompanyID int
		Name      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			"bad company id",
			fields{
				Name: "valid",
			},
			errs.ErrCompanyIDNotEmpty,
		},
		{
			"empty position name",
			fields{
				CompanyID: 1,
				Name:      "",
			},
			errs.ErrPositionNameNotEmpty,
		},
		{
			"too long position",
			fields{
				CompanyID: 1,
				Name:      randomseq.RandomString(maxPositionNameL + 1),
			},
			errs.ErrInvalidPositionName,
		},
		{
			"max length position",
			fields{
				CompanyID: 1,
				Name:      randomseq.RandomString(256),
			},
			nil,
		},
		{
			"position name conatains *",
			fields{
				CompanyID: 1,
				Name:      "*position",
			},
			errs.ErrInvalidPositionName,
		},
		{
			"position name conatains #",
			fields{
				CompanyID: 1,
				Name:      "#position",
			},
			errs.ErrInvalidPositionName,
		},
		{
			"punctuation position",
			fields{
				CompanyID: 1,
				Name:      "position,some",
			},
			nil,
		},
		{
			"space position",
			fields{
				CompanyID: 1,
				Name:      "какая-то должность",
			},
			nil,
		},
		{
			"cyrillic position",
			fields{
				CompanyID: 1,
				Name:      "стажер",
			},
			nil,
		},
		{
			"latin position",
			fields{
				CompanyID: 1,
				Name:      "validposition",
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PositionSet{
				CompanyID: tt.fields.CompanyID,
				Name:      tt.fields.Name,
			}

			err := p.Validation()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
