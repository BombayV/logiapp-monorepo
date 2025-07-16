package database

import (
	"bombayv/logiapp-monorepo/logi_api/internal/core/user"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

// UserRepository handles the database operations for users.
type UserRepository struct {
	db     *DB
	helper *DatabaseHelper
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		db:     db,
		helper: NewDatabaseHelper(db.Pool),
	}
}

// Save saves a user and their associated data to the database in a transaction.
func (r *UserRepository) Save(ctx context.Context, u *user.User, ud *user.UserData) error {
	return r.helper.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		userQuery := `
			INSERT INTO users (user_id, email, password_hash, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err := tx.Exec(ctx, userQuery, u.UserID, u.Email, u.PasswordHash, u.Role, u.CreatedAt, u.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to save user: %w", err)
		}

		userDataQuery := `
			INSERT INTO users_data (user_id, first_name, last_name, phone_number, last_connection, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.Exec(ctx, userDataQuery, ud.UserID, ud.FirstName, ud.LastName, ud.PhoneNumber, ud.LastConnection, ud.CreatedAt, ud.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to save user data: %w", err)
		}

		return nil
	})
}

// FindByEmail retrieves a user by their email address.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT user_id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	u := &user.User{}
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return u, nil
}

// FindByID retrieves a user and their data by user ID.
func (r *UserRepository) FindByID(ctx context.Context, userID string) (*user.User, *user.UserData, error) {
	var u *user.User
	var ud *user.UserData
	var err error

	err = r.helper.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Get user info
		userQuery := `
			SELECT user_id, email, password_hash, role, created_at, updated_at
			FROM users
			WHERE user_id = $1
		`
		u = &user.User{}
		err = tx.QueryRow(ctx, userQuery, userID).Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}

		// Get user data
		userDataQuery := `
			SELECT user_id, first_name, last_name, phone_number, last_connection, created_at, updated_at
			FROM users_data
			WHERE user_id = $1
		`
		ud = &user.UserData{}
		err = tx.QueryRow(ctx, userDataQuery, userID).Scan(&ud.UserID, &ud.FirstName, &ud.LastName, &ud.PhoneNumber, &ud.LastConnection, &ud.CreatedAt, &ud.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to find user data: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return u, ud, nil
}

// Update updates an existing user
func (r *UserRepository) Update(ctx context.Context, u *user.User, ud *user.UserData) error {
	return r.helper.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Update user
		userQuery := `
			UPDATE users 
			SET email = $2, password_hash = $3, role = $4, updated_at = $5
			WHERE user_id = $1
		`
		_, err := tx.Exec(ctx, userQuery, u.UserID, u.Email, u.PasswordHash, u.Role, u.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}

		// Update user data
		userDataQuery := `
			UPDATE users_data 
			SET first_name = $2, last_name = $3, phone_number = $4, last_connection = $5, updated_at = $6
			WHERE user_id = $1
		`
		_, err = tx.Exec(ctx, userDataQuery, ud.UserID, ud.FirstName, ud.LastName, ud.PhoneNumber, ud.LastConnection, ud.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to update user data: %w", err)
		}

		return nil
	})
}

// UpdateLastConnection updates the last_connection timestamp for a user
func (r *UserRepository) UpdateLastConnection(ctx context.Context, userID string) error {
	query := `
		UPDATE users_data 
		SET last_connection = NOW(), updated_at = NOW()
		WHERE user_id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last connection: %w", err)
	}
	return nil
}

// Delete deletes a user and their data (hard delete)
func (r *UserRepository) Delete(ctx context.Context, userID string) error {
	return r.helper.WithTransaction(ctx, func(ctx context.Context, tx pgx.Tx) error {
		// Delete user data first (foreign key constraint)
		_, err := tx.Exec(ctx, "DELETE FROM users_data WHERE user_id = $1", userID)
		if err != nil {
			return fmt.Errorf("failed to delete user data: %w", err)
		}

		// Delete user
		_, err = tx.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		return nil
	})
}

// Exists checks if a user exists by email
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.helper.Exists(ctx, "users", "email", email)
}

// Exists checks if a user exists by ID
func (r *UserRepository) Exists(ctx context.Context, userID string) (bool, error) {
	return r.helper.Exists(ctx, "users", "user_id", userID)
}

// FindAll retrieves all users (with pagination to avoid loading too much data)
func (r *UserRepository) FindAll(ctx context.Context, limit, offset int) ([]*user.User, int, error) {
	// Get total count
	totalCount, err := r.helper.Count(ctx, "users", "1", 1) // Count all records
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT user_id, email, password_hash, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usersList []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		usersList = append(usersList, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return usersList, totalCount, nil
}

// FindByRole retrieves users by their role
func (r *UserRepository) FindByRole(ctx context.Context, role string) ([]*user.User, error) {
	query := `
		SELECT user_id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE role = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersList []*user.User
	for rows.Next() {
		var u user.User
		if err := rows.Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		usersList = append(usersList, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return usersList, nil
}

// FindByRoleWithData retrieves users by their role with user data
func (r *UserRepository) FindByRoleWithData(ctx context.Context, role string) ([]*user.UserWithData, error) {
	query := `
		SELECT 
			u.user_id, u.email, u.password_hash, u.role, u.created_at, u.updated_at,
			ud.first_name, ud.last_name, ud.phone_number, ud.last_connection, ud.created_at as ud_created_at, ud.updated_at as ud_updated_at
		FROM users u
		LEFT JOIN users_data ud ON u.user_id = ud.user_id
		WHERE u.role = $1
		ORDER BY u.created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usersList []*user.UserWithData
	for rows.Next() {
		var u user.User
		var ud user.UserData
		var firstName, lastName, phoneNumber *string
		var lastConnection, udCreatedAt, udUpdatedAt *time.Time

		if err := rows.Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
			&firstName, &lastName, &phoneNumber, &lastConnection, &udCreatedAt, &udUpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		// Create UserWithData struct
		userWithData := &user.UserWithData{
			User: &u,
		}

		// Only populate UserData if it exists (LEFT JOIN might return nulls)
		if firstName != nil && lastName != nil && phoneNumber != nil {
			ud.UserID = u.UserID
			ud.FirstName = *firstName
			ud.LastName = *lastName
			ud.PhoneNumber = *phoneNumber
			if lastConnection != nil {
				ud.LastConnection = *lastConnection
			}
			if udCreatedAt != nil {
				ud.CreatedAt = *udCreatedAt
			}
			if udUpdatedAt != nil {
				ud.UpdatedAt = *udUpdatedAt
			}
			userWithData.UserData = &ud
		}

		usersList = append(usersList, userWithData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return usersList, nil
}

// FindAllWithData retrieves all users with their user data (with pagination)
func (r *UserRepository) FindAllWithData(ctx context.Context, limit, offset int) ([]*user.UserWithData, int, error) {
	// Get total count
	totalCount, err := r.helper.Count(ctx, "users", "1", 1) // Count all records
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results with JOIN
	query := `
		SELECT 
			u.user_id, u.email, u.password_hash, u.role, u.created_at, u.updated_at,
			ud.first_name, ud.last_name, ud.phone_number, ud.last_connection, ud.created_at as ud_created_at, ud.updated_at as ud_updated_at
		FROM users u
		LEFT JOIN users_data ud ON u.user_id = ud.user_id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usersList []*user.UserWithData
	for rows.Next() {
		var u user.User
		var ud user.UserData
		var firstName, lastName, phoneNumber *string
		var lastConnection, udCreatedAt, udUpdatedAt *time.Time

		if err := rows.Scan(&u.UserID, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt,
			&firstName, &lastName, &phoneNumber, &lastConnection, &udCreatedAt, &udUpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		// Create UserWithData struct
		userWithData := &user.UserWithData{
			User: &u,
		}

		// Only populate UserData if it exists (LEFT JOIN might return nulls)
		if firstName != nil && lastName != nil && phoneNumber != nil {
			ud.UserID = u.UserID
			ud.FirstName = *firstName
			ud.LastName = *lastName
			ud.PhoneNumber = *phoneNumber
			if lastConnection != nil {
				ud.LastConnection = *lastConnection
			}
			if udCreatedAt != nil {
				ud.CreatedAt = *udCreatedAt
			}
			if udUpdatedAt != nil {
				ud.UpdatedAt = *udUpdatedAt
			}
			userWithData.UserData = &ud
		}

		usersList = append(usersList, userWithData)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return usersList, totalCount, nil
}

// UpdateLocation updates or inserts a user's location
func (r *UserRepository) UpdateLocation(ctx context.Context, userID string, latitude, longitude float64) error {
	query := `
		INSERT INTO users_locations (user_id, location, updated_at)
		VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326), NOW())
		ON CONFLICT (user_id)
		DO UPDATE SET location = ST_SetSRID(ST_MakePoint($2, $3), 4326), updated_at = NOW()
	`
	_, err := r.db.Pool.Exec(ctx, query, userID, longitude, latitude)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}
	return nil
}

// GetLocation retrieves a user's location
func (r *UserRepository) GetLocation(ctx context.Context, userID string) (*user.UserLocation, error) {
	query := `
		SELECT user_id, ST_X(location) as longitude, ST_Y(location) as latitude, updated_at
		FROM users_locations
		WHERE user_id = $1
	`
	var userLocation user.UserLocation
	var latitude, longitude float64

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&userLocation.UserID,
		&longitude,
		&latitude,
		&userLocation.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("location not found for user")
		}
		return nil, fmt.Errorf("failed to get location: %w", err)
	}

	// Set the computed fields
	userLocation.Latitude = latitude
	userLocation.Longitude = longitude
	userLocation.Location = fmt.Sprintf("POINT(%f %f)", longitude, latitude)

	return &userLocation, nil
}

// GetActiveDriversWithLocations retrieves all drivers who have been active in the last 10 minutes with their locations
func (r *UserRepository) GetActiveDriversWithLocations(ctx context.Context) ([]*user.DriverLocation, error) {
	query := `
		SELECT 
			u.user_id,
			u.email,
			ud.first_name,
			ud.last_name,
			ud.phone_number,
			ST_X(ul.location) as longitude,
			ST_Y(ul.location) as latitude,
			ud.last_connection,
			ul.updated_at as location_updated_at
		FROM users u
		INNER JOIN users_data ud ON u.user_id = ud.user_id
		INNER JOIN users_locations ul ON u.user_id = ul.user_id
		WHERE u.role = 'driver' 
		AND ud.last_connection > NOW() - INTERVAL '10 minutes'
		ORDER BY ud.last_connection DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query active drivers: %w", err)
	}
	defer rows.Close()

	var drivers []*user.DriverLocation
	for rows.Next() {
		var driver user.DriverLocation
		var latitude, longitude float64

		err := rows.Scan(
			&driver.UserID,
			&driver.Email,
			&driver.FirstName,
			&driver.LastName,
			&driver.PhoneNumber,
			&longitude,
			&latitude,
			&driver.LastConnection,
			&driver.LocationUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}

		// Set the computed fields
		driver.Latitude = latitude
		driver.Longitude = longitude

		drivers = append(drivers, &driver)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return drivers, nil
}

// GetAllDrivers retrieves all drivers ordered by last connection
func (r *UserRepository) GetAllDrivers(ctx context.Context) ([]*user.Driver, error) {
	query := `
		SELECT 
			u.user_id,
			u.email,
			ud.first_name,
			ud.last_name,
			ud.phone_number,
			u.role,
			ud.last_connection,
			u.created_at,
			u.updated_at
		FROM users u
		INNER JOIN users_data ud ON u.user_id = ud.user_id
		WHERE u.role = 'driver'
		ORDER BY ud.last_connection DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []*user.Driver
	for rows.Next() {
		var driver user.Driver

		err := rows.Scan(
			&driver.UserID,
			&driver.Email,
			&driver.FirstName,
			&driver.LastName,
			&driver.PhoneNumber,
			&driver.Role,
			&driver.LastConnection,
			&driver.CreatedAt,
			&driver.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}

		drivers = append(drivers, &driver)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return drivers, nil
}
