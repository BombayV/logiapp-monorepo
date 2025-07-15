package user

import (
	"bombayv/logiapp-monorepo/logi_api/internal/storage/cache"
	"bombayv/logiapp-monorepo/logi_api/internal/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("missing or invalid credentials")
	ErrPasswordTooWeak    = errors.New("password must be at least 8 characters long")
	ErrForbidden          = errors.New("forbidden: Insufficient permissions")
	ErrEmailExists        = errors.New("email address is already in use")
	ErrCouldNotSaveUser   = errors.New("could not save user to the database")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

// Service provides user-related operations.
type Service struct {
	repo  Repository
	cache *cache.Cache
}

// NewService creates a new user service.
func NewService(repo Repository, cache *cache.Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

// CreateUser handles the business logic for creating a new user.
func (s *Service) CreateUser(ctx context.Context, email, password, firstName, lastName, phone, role string) (*User, error) {
	// 2. Validate role
	if role != "sales" && role != "driver" {
		return nil, errors.New("invalid role provided: must be 'sales', or 'driver'")
	}

	// 3. Validate email and password strength.
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	if len(password) < 8 {
		return nil, ErrPasswordTooWeak
	}

	// 4. Hash the password.
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// 5. Create User and UserData structs.
	userID := uuid.New().String()
	now := time.Now()
	user := &User{
		UserID:       userID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	userData := &UserData{
		UserID:         userID,
		FirstName:      firstName,
		LastName:       lastName,
		PhoneNumber:    phone,
		LastConnection: now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// 6. Call the repository to save the user to the database.
	if err := s.repo.Save(ctx, user, userData); err != nil {
		return nil, err
	}

	// 7. Return the created user.
	return user, nil
}

// Login authenticates a user and returns a JWT token if successful.
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	ok, err := utils.CheckPasswordHash(password, user.PasswordHash)
	if err != nil || !ok {
		return "", ErrInvalidCredentials
	}

	return utils.GenerateJWT(user.UserID, user.Role)
}

// Logout revokes a user's JWT token.
func (s *Service) Logout(ctx context.Context, tokenString string) error {
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		return ErrInvalidToken
	}

	// Add the token's JTI to the revocation list in Redis.
	// The token is stored until its original expiration time.
	expiresAt := claims.ExpiresAt.Time
	if time.Now().After(expiresAt) {
		return ErrInvalidToken
	}

	err = s.cache.AddRevokedToken(ctx, claims.ID, time.Until(expiresAt))
	if err != nil {
		return err
	}

	// Update the user's last_connection timestamp
	if err := s.repo.UpdateLastConnection(ctx, claims.Subject); err != nil {
		// Log the error but don't fail the logout process
		// The token is already revoked, so the user is effectively logged out
		return fmt.Errorf("failed to update last connection: %w", err)
	}

	return nil
}

// GetUserProfile retrieves a user's profile information by user ID.
func (s *Service) GetUserProfile(ctx context.Context, userID string) (*User, *UserData, error) {
	user, userData, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	return user, userData, nil
}

// GetAllUsers retrieves all users with pagination (admin function).
func (s *Service) GetAllUsers(ctx context.Context, limit, offset int) ([]*User, int, error) {
	users, total, err := s.repo.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetAllUsersWithData retrieves all users with their profile data (admin function).
func (s *Service) GetAllUsersWithData(ctx context.Context, limit, offset int) ([]*UserWithData, int, error) {
	users, total, err := s.repo.FindAllWithData(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUsersByRole retrieves users filtered by role.
func (s *Service) GetUsersByRole(ctx context.Context, role string) ([]*User, error) {
	// Validate role
	if role != "admin" && role != "driver" && role != "sales" {
		return nil, errors.New("invalid role provided")
	}

	users, err := s.repo.FindByRole(ctx, role)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersByRoleWithData retrieves users filtered by role with their profile data.
func (s *Service) GetUsersByRoleWithData(ctx context.Context, role string) ([]*UserWithData, error) {
	// Validate role
	if role != "admin" && role != "driver" && role != "sales" {
		return nil, errors.New("invalid role provided")
	}

	users, err := s.repo.FindByRoleWithData(ctx, role)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserProfile updates a user's profile information.
func (s *Service) UpdateUserProfile(ctx context.Context, userID string, firstName, lastName, phone, role *string) (*User, *UserData, error) {
	// Get current user data
	user, userData, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	// Update fields if provided
	if firstName != nil {
		userData.FirstName = *firstName
	}
	if lastName != nil {
		userData.LastName = *lastName
	}
	if phone != nil {
		userData.PhoneNumber = *phone
	}
	if role != nil {
		// Validate role
		if *role != "admin" && *role != "driver" && *role != "sales" {
			return nil, nil, errors.New("invalid role provided")
		}
		user.Role = *role
	}

	userData.UpdatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save updates
	err = s.repo.Update(ctx, user, userData)
	if err != nil {
		return nil, nil, err
	}

	return user, userData, nil
}

// DeleteUser deletes a user (admin function).
func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	// Check if user exists
	exists, err := s.repo.Exists(ctx, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user not found")
	}

	err = s.repo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateLocation updates a user's location
func (s *Service) UpdateLocation(ctx context.Context, userID string, latitude, longitude float64) error {
	// Validate latitude and longitude ranges
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90")
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180")
	}

	return s.repo.UpdateLocation(ctx, userID, latitude, longitude)
}

// GetLocation retrieves a user's location
func (s *Service) GetLocation(ctx context.Context, userID string) (*UserLocation, error) {
	return s.repo.GetLocation(ctx, userID)
}

// GetActiveDriversWithLocations retrieves all drivers who have been active in the last 10 minutes with their locations
func (s *Service) GetActiveDriversWithLocations(ctx context.Context) ([]*DriverLocation, error) {
	return s.repo.GetActiveDriversWithLocations(ctx)
}

// ResetPassword changes a user's password after verifying the current password
func (s *Service) ResetPassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Validate new password strength
	if len(newPassword) < 8 {
		return ErrPasswordTooWeak
	}

	// Get user by ID to verify current password
	user, userData, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return ErrInvalidCredentials
	}

	// Verify current password
	ok, err := utils.CheckPasswordHash(currentPassword, user.PasswordHash)
	if err != nil || !ok {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password in database
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	// Update user with new password
	err = s.repo.Update(ctx, user, userData)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
