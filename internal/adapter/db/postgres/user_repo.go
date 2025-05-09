package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"ip_detector/internal/domain/model"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Save(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, email, ip, country, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.IP,
		user.Country,
		user.PasswordHash,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

func (r *PostgresUserRepo) GetAll(ctx context.Context) ([]*model.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, email, ip, country FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IP, &u.Country); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}

	return users, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, ip, country FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email, &user.IP, &user.Country)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (r *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, ip, country, password_hash FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.IP, &user.Country, &user.PasswordHash)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}
