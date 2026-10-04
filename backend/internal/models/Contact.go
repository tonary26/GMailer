package models

import "time"

type Contact struct {
	ID 			int 		`json:"id" db:"id"`
	Email		string		`json:"email" db:"email"`
	Name 		string 		`json:"name" db:"name"` 
	CreatedAt  time.Time	`json:"created_at" db:"created_at"`
}

type ContactCreateInput struct {
	Email		string		`json:"email" validate:"required,email"`
	Name 		string 		`json:"name" validate:"required"` 
}
