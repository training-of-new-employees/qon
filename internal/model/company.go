package model

type (
	Company struct {
		ID         int    `db:"id" json:"id"`
		IsActive   bool   `db:"active" json:"active"`
		IsArchived bool   `db:"archived" json:"archived"`
		Name       string `db:"name" json:"name"`
		CreatedAt  string `db:"created_at" json:"created_at"`
		UpdatedAt  string `db:"updated_at" json:"updated_at"`
	}

	CompanyEdit struct {
		ID         int     `json:"-" db:"id"`
		IsActive   *bool   `json:"active" db:"active"`
		IsArchived *bool   `json:"archived" db:"archived"`
		Name       *string `json:"name,omitempty" db:"name"`
	}
)
