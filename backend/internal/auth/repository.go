package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(email, username, passwordHash string) (User, error)
	GetByID(id int64) (User, bool, error)
	GetByIdentifier(identifier string) (userRecord, bool, error)
	ClaimOrphanedData(userID int64) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(email, username, passwordHash string) (User, error) {
	var user User

	query := `
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, created_at
	`

	err := r.db.QueryRow(context.Background(), query, email, username, passwordHash).
		Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return User{}, ErrEmailAlreadyInUse
			case "users_username_key":
				return User{}, ErrUsernameAlreadyInUse
			default:
				message := strings.ToLower(pgErr.Message)
				if strings.Contains(message, "email") {
					return User{}, ErrEmailAlreadyInUse
				}
				if strings.Contains(message, "username") {
					return User{}, ErrUsernameAlreadyInUse
				}
			}
		}

		return User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) GetByID(id int64) (User, bool, error) {
	var user User

	query := `
		SELECT id, email, username, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(context.Background(), query, id).
		Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, false, nil
		}

		return User{}, false, err
	}

	return user, true, nil
}

func (r *PostgresRepository) GetByIdentifier(identifier string) (userRecord, bool, error) {
	var record userRecord

	query := `
		SELECT id, email, username, password_hash, created_at
		FROM users
		WHERE email = $1 OR username = $1
	`

	err := r.db.QueryRow(context.Background(), query, identifier).
		Scan(&record.ID, &record.Email, &record.Username, &record.PasswordHash, &record.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userRecord{}, false, nil
		}

		return userRecord{}, false, err
	}

	return record, true, nil
}

func (r *PostgresRepository) ClaimOrphanedData(userID int64) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	if _, err := tx.Exec(context.Background(), `UPDATE accounts SET user_id = $1 WHERE user_id IS NULL`, userID); err != nil {
		return err
	}

	if _, err := tx.Exec(context.Background(), `UPDATE transactions SET user_id = $1 WHERE user_id IS NULL`, userID); err != nil {
		return err
	}

	return tx.Commit(context.Background())
}
