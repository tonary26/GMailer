package database

import (
	"context"
	"errors"
	"gmailer/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthStore struct {
	db *pgxpool.Pool
}

func NewAuthStore(db *pgxpool.Pool) *AuthStore {
	return &AuthStore{db: db}
}

func (s *AuthStore) Register(ctx context.Context, input models.UserCreateInput) (*models.User, error) {
	var user models.User

	query := `
		INSERT INTO users 
		(email, password_hash) 
		VALUES ($1, $2) 
		RETURNING id, email, password_hash, created_at;
	`

	err := s.db.QueryRow(ctx, query, input.Email, input.Password).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	query := `
		SELECT id, email, password_hash, created_at
		FROM users 
		WHERE email = $1;
	`

	err := s.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
} 