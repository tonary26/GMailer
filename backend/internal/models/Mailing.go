package models

import "time"

type Mailing struct {
	ID           int        `json:"id" db:"id"`
	Title        string     `json:"title" db:"title"`
	Subject      string     `json:"subject" db:"subject"`
	BodyTemplate string     `json:"body_template" db:"body_template"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	StartedAt    *time.Time `json:"started_at" db:"started_at"`
}

type CreateMailingInput struct {
	Title        string `json:"title" db:"title"`
	Subject      string `json:"subject" db:"subject"`
	BodyTemplate string `json:"body_template" db:"body_template"`
}
