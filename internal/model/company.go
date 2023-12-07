package model

type (
	Company struct {
		ID        int    `db:"id" json:"id"`
		Name      string `db:"name" json:"name"`
		IsActive  bool   `db:"active" json:"active"`
		IsDeleted bool   `db:"is_deleted" json:"is_deleted"`
		CreatedAt string `db:"created_at" json:"created_at"`
		UpdatedAt string `db:"updated_at" json:"updated_at"`
	}
)
