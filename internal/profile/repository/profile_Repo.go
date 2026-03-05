package repository

import (
	"cloud-kitchen/internal/auth/model"
	profilemodel "cloud-kitchen/internal/profile/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// Implement methods for profile repository (e.g., GetProfile, UpdateProfile, DeleteProfile)
func (r *ProfileRepository) GetProfile(ctx context.Context, email string) (*profilemodel.Profile, error) {
	query := `
	SELECT id, name, email,provider, mobile_number, profile_picture
	FROM users
	WHERE email = $1
	`
	row := r.db.QueryRow(ctx, query, email)
	var profile profilemodel.Profile
	err := row.Scan(
		&profile.ID,
		&profile.Name,
		&profile.Email,
		&profile.Provider,
		&profile.MobileNumber,
		&profile.ProfilePicture,
	)
	if err != nil {
		return nil, err
	}

	queryAddresses := `
	SELECT id, user_id, label, street, city, state, zip_code,latitude, longitude, is_default
	FROM addresses
	WHERE user_id =  $1
	`
	rows, err := r.db.Query(ctx, queryAddresses, profile.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.AddressModel
	for rows.Next() {
		var address model.AddressModel
		err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.Label,
			&address.Street,
			&address.City,
			&address.State,
			&address.ZipCode,
			&address.Latitude,
			&address.Longitude,
			&address.IsDefault,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	profile.Addresses = addresses
	return &profile, nil
}

func (r *ProfileRepository) UpdateUserProfile(ctx context.Context, userID string, req *profilemodel.UpdateProfileRequest) error {

	query := `
	UPDATE users
	SET name=$1,
	    mobile_number=$2,
	    profile_picture=$3
	WHERE id=$4
	`

	_, err := r.db.Exec(
		ctx,
		query,
		req.Name,
		req.MobileNumber,
		req.ProfilePicture,
		userID,
	)

	return err
}

func (r *ProfileRepository) UpdateAddress(ctx context.Context, addr *model.AddressModel) error {

	query := `
	UPDATE addresses
	SET label=$1,
	    street=$2,
	    city=$3,
	    state=$4,
	    zip_code=$5,
	    latitude=$6,
	    longitude=$7,
	    is_default=$8
	WHERE id=$9 AND user_id=$10
	`

	_, err := r.db.Exec(
		ctx,
		query,
		addr.Label,
		addr.Street,
		addr.City,
		addr.State,
		addr.ZipCode,
		addr.Latitude,
		addr.Longitude,
		addr.IsDefault,
		addr.ID,
		addr.UserID,
	)

	return err
}

func (r *ProfileRepository) CreateAddress(ctx context.Context, addr *model.AddressModel) error {

	query := `
	INSERT INTO addresses
	(id, user_id, label, street, city, state, zip_code, latitude, longitude, is_default)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		addr.ID,
		addr.UserID,
		addr.Label,
		addr.Street,
		addr.City,
		addr.State,
		addr.ZipCode,
		addr.Latitude,
		addr.Longitude,
		addr.IsDefault,
	)

	return err
}

func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID string) error {

	query := `
	DELETE FROM users
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, userID)

	return err
}
