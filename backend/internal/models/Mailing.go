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

type PendingMessage struct {
	ID           int    `db:"id"`
	MailingID    int    `db:"mailing_id"`
	ContactID    int    `db:"contact_id"`
	Email        string `db:"email"`
	Name         string `db:"name"`
	Subject      string `db:"subject"`
	BodyTemplate string `db:"body_template"`
}

type MailingProgress struct {
	MailingID int    `json:"mailing_id" db:"mailing_id"`
	Status    string `json:"status" db:"status"`
	Total     int    `json:"total" db:"total"`
	Pending   int    `json:"pending" db:"pending"`
	Sending   int    `json:"sending" db:"sending"`
	Sent      int    `json:"sent" db:"sent"`
	Failed    int    `json:"failed" db:"failed"`
}

type CreateMailingInput struct {
	Title        string `json:"title" db:"title"`
	Subject      string `json:"subject" db:"subject"`
	BodyTemplate string `json:"body_template" db:"body_template"`
}
