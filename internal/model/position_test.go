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
			"Bad company ID",
			fields{
				Name: "valid",
			},
			errs.ErrCompanyIDNotEmpty,
		},
		{
			"Short position",
			fields{
				CompanyID: 1,
				Name:      "s",
			},
			errs.ErrInvalidPositionName,
		},
		{
			"Long position",
			fields{
				CompanyID: 1,
				Name:      randomseq.RandomString(257),
			},
			errs.ErrInvalidPositionName,
		},
		{
			"Correct length position",
			fields{
				CompanyID: 1,
				Name:      randomseq.RandomString(256),
			},
			nil,
		},
		{
			"Bad symbols position",
			fields{
				CompanyID: 1,
				Name:      "*position",
			},
			errs.ErrInvalidPositionName,
		},
		{
			"Bad symbols position",
			fields{
				CompanyID: 1,
				Name:      "#position",
			},
			errs.ErrInvalidPositionName,
		},
		{
			"Punctuation position",
			fields{
				CompanyID: 1,
				Name:      "position,some",
			},
			errs.ErrInvalidPositionName,
		},
		{
			"Space position",
			fields{
				CompanyID: 1,
				Name:      "какая-то должность",
			},
			nil,
		},
		{
			"Cyrillic position",
			fields{
				CompanyID: 1,
				Name:      "стажер",
			},
			nil,
		},
		{
			"Latin position",
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
