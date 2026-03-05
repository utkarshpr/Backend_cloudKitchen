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

func (r *ProfileRepository) UpdateProfile(ctx context.Context, profile *profilemodel.Profile) error {
	// Implement logic to update profile in the database using r.db
	return nil
}

func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID int) error {
	// Implement logic to delete profile from the database using r.db
	return nil
}
