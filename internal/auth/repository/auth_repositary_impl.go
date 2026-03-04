package repository

import (
	"context"

	"cloud-kitchen/internal/auth/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresAuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) AuthRepository {
	return &PostgresAuthRepository{
		db: db,
	}
}

func (r *PostgresAuthRepository) CreateUser(user *model.User) error {

	query := `
	INSERT INTO users (id, name, email, password, provider)
	VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Provider,
	)

	return err
}

func (r *PostgresAuthRepository) GetUserByEmail(email string) (*model.User, error) {

	query := `
	SELECT id, name, email, password, provider
	FROM users
	WHERE email=$1
	`

	row := r.db.QueryRow(context.Background(), query, email)

	var user model.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
	)

	if err != nil {
		return nil, nil
	}

	return &user, nil
}
