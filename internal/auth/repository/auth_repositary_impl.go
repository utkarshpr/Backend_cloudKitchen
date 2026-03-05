package repository

import (
	"context"
	"errors"

	"cloud-kitchen/internal/auth/model"
	"cloud-kitchen/pkg/util"

	"github.com/jackc/pgx/v5"
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
INSERT INTO users
(id, name, email, password, provider, mobile_number, profile_picture, created_at)
VALUES ($1,$2,$3,$4,$5,$6,$7,NOW())
`

	_, err := r.db.Exec(
	ctx,
	query,
	user.ID,
	user.Name,
	user.Email,
	user.Password,
	user.Provider,
	user.MobileNumber,
	user.ProfilePicture,
)

	if err != nil {
		util.Error(ctx, "repository.CreateUser error=%v", err)
		return err
	}

	return nil
}

func (r *PostgresAuthRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {

	util.Info(ctx, "repository.GetUserByEmail email=%s", email)

	query := `
	SELECT id, name, email, password, provider, mobile_number, profile_picture, created_at
	FROM users
	WHERE email=$1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
		&user.MobileNumber,
		&user.ProfilePicture,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			util.Info(ctx, "repository.GetUserByEmail user not found email=%s", email)
			return nil, nil
		}

		util.Error(ctx, "repository.GetUserByEmail db error=%v", err)
		return nil, err
	}

	return &user, nil
}

func (r *PostgresAuthRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {

	util.Info(ctx, "repository.GetUserByID id=%s", id)

	query := `
	SELECT id, name, email, password, provider, created_at
	FROM users
	WHERE id=$1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Provider,
		&user.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			util.Info(ctx, "repository.GetUserByID user not found id=%s", id)
			return nil, nil
		}

		util.Error(ctx, "repository.GetUserByID db error=%v", err)
		return nil, err
	}

	return &user, nil
}

func (r *PostgresAuthRepository) CreateAddress(ctx context.Context, address *model.AddressModel) error {

	query := `
	INSERT INTO addresses
	(id, user_id, label, street, city, state, zip_code, latitude, longitude, is_default)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		address.ID,
		address.UserID,
		address.Label,
		address.Street,
		address.City,
		address.State,
		address.ZipCode,
		address.Latitude,
		address.Longitude,
		address.IsDefault,
	)

	return err
}

func (r *PostgresAuthRepository) GetAddressesByUserID(ctx context.Context, userID string) ([]model.Address, error) {

	query := `
	SELECT label, street, city, state, zip_code, latitude, longitude, is_default
	FROM addresses
	WHERE user_id=$1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address

	for rows.Next() {

		var addr model.Address

		err := rows.Scan(
			&addr.Label,
			&addr.Street,
			&addr.City,
			&addr.State,
			&addr.ZipCode,
			&addr.Latitude,
			&addr.Longitude,
			&addr.IsDefault,
		)

		if err != nil {
			return nil, err
		}

		addresses = append(addresses, addr)
	}

	return addresses, nil
}