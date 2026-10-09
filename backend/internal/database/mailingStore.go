package database

import (
	"context"
	"errors"
	"gmailer/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNoContacts      = errors.New("mailing requires at least one contact")
	ErrMailingNotDraft = errors.New("mailing is not a draft")
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

	var status string
	if err := tx.QueryRow(ctx, `SELECT status FROM mailings WHERE id = $1 FOR UPDATE`, mailingId).Scan(&status); err != nil {
		return err
	}
	if status != "draft" {
		return ErrMailingNotDraft
	}

	var contactsCount int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM contacts`).Scan(&contactsCount); err != nil {
		return err
	}
	if contactsCount == 0 {
		return ErrNoContacts
	}

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

func (s *MailingStore) FetchPending(ctx context.Context, limit int) ([]models.PendingMessage, error) {
	query := `
		WITH picked AS (
			SELECT id FROM mailing_messages
			WHERE status = 'pending' AND next_attempt_at <= now()
			ORDER BY id
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE mailing_messages m
		SET status = 'sending', locked_at = now(), attempts = m.attempts + 1
		FROM picked, contacts c, mailings ml
		WHERE m.id = picked.id AND m.contact_id = c.id AND m.mailing_id = ml.id
		RETURNING m.id, m.mailing_id, m.contact_id, c.email, c.name, ml.subject, ml.body_template
	`

	rows, err := s.db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.PendingMessage])
}

func (s *MailingStore) MarkSent(ctx context.Context, id int) error {
	query := `
		UPDATE mailing_messages SET status = 'sent', sent_at = now() WHERE id = $1
	`
	_, err := s.db.Exec(ctx, query, id)
	return err
}

func (s *MailingStore) MarkFailed(ctx context.Context, id int, maxAttempts int, backoff time.Duration) error {
	query := `
		UPDATE mailing_messages
		SET
			status = CASE WHEN attempts >= $2 THEN 'failed' ELSE 'pending' END,
			next_attempt_at = now() + $3
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, id, maxAttempts, backoff)
	return err
}

func (s *MailingStore) MarkDoneIfComplete(ctx context.Context, mailingID int) error {
	query := `
		UPDATE mailings m
		SET status = 'done'
		WHERE m.id = $1
			AND m.status = 'running'
			AND EXISTS (SELECT 1 FROM mailing_messages WHERE mailing_id = m.id)
			AND NOT EXISTS (
				SELECT 1 FROM mailing_messages
				WHERE mailing_id = m.id AND status IN ('pending', 'sending')
			)
	`
	_, err := s.db.Exec(ctx, query, mailingID)
	return err
}

func (s *MailingStore) ListProgress(ctx context.Context) ([]models.MailingProgress, error) {
	query := `
		SELECT
			m.id AS mailing_id,
			m.status,
			count(mm.id)::int AS total,
			count(mm.id) FILTER (WHERE mm.status = 'pending')::int AS pending,
			count(mm.id) FILTER (WHERE mm.status = 'sending')::int AS sending,
			count(mm.id) FILTER (WHERE mm.status = 'sent')::int AS sent,
			count(mm.id) FILTER (WHERE mm.status = 'failed')::int AS failed
		FROM mailings m
		LEFT JOIN mailing_messages mm ON mm.mailing_id = m.id
		GROUP BY m.id, m.status
		ORDER BY m.id DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[models.MailingProgress])
}
