package database

import (
	"context"
	"gmailer/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactStore struct {
	db *pgxpool.Pool
}

func NewContactStore(db *pgxpool.Pool) *ContactStore {
	return &ContactStore{db: db}
}

func (s *ContactStore) CreateContact(ctx context.Context, input models.ContactCreateInput) (*models.Contact, error) {
	var contact models.Contact

	query := `
		INSERT INTO contacts (email, name) VALUES ($1, $2)
		RETURNING id, email, name, created_at
	`

	err := s.db.QueryRow(ctx, query, input.Email, input.Name).Scan(
		&contact.ID,
		&contact.Email,
		&contact.Name,
		&contact.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &contact, nil
}

func (s *ContactStore) ListContact(ctx context.Context) (*[]models.Contact, error) {
	query := `
		SELECT id, email, name, created_at FROM contacts ORDER BY id DESC
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Contact])
	if err != nil {
		return nil, err
	}

	return &contacts, nil
}