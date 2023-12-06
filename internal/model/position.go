package model

import "time"

type (
	Position struct {
		ID        int       `db:"id" json:"id"`
		CompanyID int       `db:"company_id" json:"company_id"`
		Name      string    `db:"name" json:"name"`
		IsActive  bool      `db:"active" json:"active"`
		CreatedAt time.Time `db:"created_at" json:"created_at"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}

	CreatePosition struct {
		CompanyID string `json:"company_id" db:"company_id"`
		Name      string `json:"name" db:"name"`
	}
)
