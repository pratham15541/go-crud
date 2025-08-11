package services

import (
	"fmt"
	"strings"

	"github.com/pratham15541/go-crud/internal/models"
	"github.com/pratham15541/go-crud/internal/repository"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req *models.CreateUserRequest) (*models.User, error) {
	// Validate business rules
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Create user
	user, err := s.userRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUsers retrieves all users with pagination
func (s *UserService) GetUsers(page, limit int) ([]*models.User, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Get users
	users, err := s.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	// Get total count
	total, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id int, req *models.UpdateUserRequest) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Validate update request
	if err := s.validateUpdateUserRequest(req); err != nil {
		return nil, err
	}

	// Check if email is being updated and already exists
	if req.Email != "" {
		existingUser, _ := s.userRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	// Update user
	user, err := s.userRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	err := s.userRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// validateCreateUserRequest validates create user request
func (s *UserService) validateCreateUserRequest(req *models.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if len(req.Name) < 2 || len(req.Name) > 100 {
		return fmt.Errorf("name must be between 2 and 100 characters")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("email is required")
	}

	if !isValidEmail(req.Email) {
		return fmt.Errorf("invalid email format")
	}

	if req.Age <= 0 || req.Age > 150 {
		return fmt.Errorf("age must be between 1 and 150")
	}

	return nil
}

// validateUpdateUserRequest validates update user request
func (s *UserService) validateUpdateUserRequest(req *models.UpdateUserRequest) error {
	if req.Name != "" {
		if len(req.Name) < 2 || len(req.Name) > 100 {
			return fmt.Errorf("name must be between 2 and 100 characters")
		}
	}

	if req.Email != "" {
		if !isValidEmail(req.Email) {
			return fmt.Errorf("invalid email format")
		}
	}

	if req.Age != 0 {
		if req.Age <= 0 || req.Age > 150 {
			return fmt.Errorf("age must be between 1 and 150")
		}
	}

	return nil
}

// isValidEmail validates email format (basic validation)
func isValidEmail(email string) bool {
	// Basic email validation
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}