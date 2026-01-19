package models

type Student struct {
	ID        int    `json:"id,omitempty" db:"id,omitempty" validate:"required"`
	FirstName string `json:"first_name,omitempty" db:"first_name,omitempty" validate:"required"`
	LastName  string `json:"last_name,omitempty" db:"last_name,omitempty" validate:"required"`
	Email     string `json:"email,omitempty" db:"email,omitempty" validate:"required"`
	Class     string `json:"class,omitempty" db:"class,omitempty" validate:"required"`
}
