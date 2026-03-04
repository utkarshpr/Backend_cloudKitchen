package repository

import (
	"context"

	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/pkg/util"

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

func (r *PostgresAuthRepository) CreateUser(ctx context.Context, user *model.User) error {
	util.Info(ctx, "repository.CreateUser user=%s email=%s", user.ID, user.Email)

	query := `
	INSERT INTO users (id, name, email, password, provider)
	VALUES ($1,$2,$3,$4,$5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Provider,
	)

	if err != nil {
		util.Error(ctx, "repository.CreateUser error=%v", err)
	}
	return err
}

func (r *PostgresAuthRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	util.Info(ctx, "repository.GetUserByEmail email=%s", email)

	query := `
	SELECT id, name, email, password, provider
	FROM users
	WHERE email=$1
	`

	row := r.db.QueryRow(ctx, query, email)

	var user model.User

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
	)

	if err != nil {
		util.Info(ctx, "repository.GetUserByEmail not found or error=%v", err)
		return nil, nil
	}

	return &user, nil
}
