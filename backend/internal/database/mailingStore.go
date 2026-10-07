package database

import (
	"context"
	"gmailer/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MailingStore struct {
	db *pgxpool.Pool
}

func NewMailingStore(db *pgxpool.Pool) *MailingStore {
	return &MailingStore{db: db}
}

func (s *MailingStore) CreateMailing(ctx context.Context, input models.CreateMailingInput) (*models.Mailing, error) {
	var mailing models.Mailing

	query := `
		INSERT INTO mailings (title, subject, body_template)
		VALUES ($1, $2, $3)
		RETURNING id, title, subject, body_template, status, created_at, started_at
	`

	err := s.db.QueryRow(ctx, query, input.Title, input.Subject, input.BodyTemplate).Scan(
		&mailing.ID,
		&mailing.Title,
		&mailing.Subject,
		&mailing.BodyTemplate,
		&mailing.Status,
		&mailing.CreatedAt,
		&mailing.StartedAt,
	)
	if err != nil {
		return nil, err
	}

	return &mailing, nil
}

func (s *MailingStore) ListMailings(ctx context.Context) (*[]models.Mailing, error) {
	query := `
		SELECT id, title, subject, body_template, status, created_at, started_at
		FROM mailings ORDER BY id DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mailing, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Mailing])
	if err != nil {
		return nil, err
	}

	return &mailing, nil
}

func (s *MailingStore) GetMailingById(ctx context.Context, id int64) (*models.Mailing, error) {
	var mailing models.Mailing

	query := `
		SELECT id, title, subject, body_template, status, created_at, started_at
		FROM mailings WHERE id = $1
	`

	err := s.db.QueryRow(ctx, query, id).Scan(
		&mailing.ID,
		&mailing.Title,
		&mailing.Subject,
		&mailing.BodyTemplate,
		&mailing.Status,
		&mailing.CreatedAt,
		&mailing.StartedAt,
	)
	if err != nil {
		return nil, err
	}

	return &mailing, nil
}

func (s *MailingStore) Start(ctx context.Context, mailingId int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO mailing_messages (mailing_id, contact_id)
		SELECT $1, id FROM contacts
		ON CONFLICT (mailing_id, contact_id) DO NOTHING
	`

	_, err = tx.Exec(ctx, query, mailingId)
	if err != nil {
		return err
	}

	query = `
		UPDATE mailings SET status = 'running', started_at = now() WHERE id = $1
	`

	_, err = tx.Exec(ctx, query, mailingId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
