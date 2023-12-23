package model

import (
	"testing"
)

func TestPositionEdit_Validation(t *testing.T) {
	type fields struct {
		CompanyID int
		Name      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Bad company ID",
			fields{
				Name: "valid",
			},
			true,
		},
		{
			"Short position",
			fields{
				CompanyID: 1,
				Name:      "s",
			},
			true,
		},
		{
			"Long position",
			fields{
				CompanyID: 1,
				Name:      "kjaskdfjajuydiuywreiqpoeripwqoeiqwumdreupoidwqueirdupweqouiruqwpmoeiurdoiquwepmraisodifpuop2iuriuoiufosudofiauodufoasufouasdoifupoisuafajdkfjaslkjdfkjaskfjdsalkjfkalsdjflajsdfkasjofapsudfaopsiufdaisjfqkjwkejlkjfasdkljfoqwpeasdfjalskdjflsajdlfkjasdklfjlaksjfdlkasjdfkasjdfhajskhfu",
			},
			true,
		},
		{
			"Bad symbols position",
			fields{
				CompanyID: 1,
				Name:      "#position",
			},
			true,
		},
		{
			"Punctuation position",
			fields{
				CompanyID: 1,
				Name:      "position,some",
			},
			true,
		},
		{
			"Space position",
			fields{
				CompanyID: 1,
				Name:      "какая-то должность",
			},
			true,
		},
		{
			"Cyrillic position",
			fields{
				CompanyID: 1,
				Name:      "стажер",
			},
			false,
		},
		{
			"Latin position",
			fields{
				CompanyID: 1,
				Name:      "validposition",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PositionEdit{
				CompanyID: tt.fields.CompanyID,
				Name:      tt.fields.Name,
			}
			if err := p.Validation(); (err != nil) != tt.wantErr {
				t.Errorf("PositionEdit.Validation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
