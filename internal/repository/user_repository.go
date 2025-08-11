package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/pratham15541/go-crud/internal/models"
)

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(req *models.CreateUserRequest) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, age) 
		VALUES ($1, $2, $3) 
		RETURNING id, name, email, age, created_at, updated_at
	`

	user := &models.User{}
	err := r.db.QueryRow(query, req.Name, req.Email, req.Age).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetAll retrieves all users with pagination
func (r *userRepository) GetAll(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

// Update updates a user
func (r *userRepository) Update(id int, req *models.UpdateUserRequest) (*models.User, error) {
	// First, get the current user
	currentUser, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.Name != "" {
		currentUser.Name = req.Name
	}
	if req.Email != "" {
		currentUser.Email = req.Email
	}
	if req.Age != 0 {
		currentUser.Age = req.Age
	}
	currentUser.UpdatedAt = time.Now()

	query := `
		UPDATE users 
		SET name = $1, email = $2, age = $3, updated_at = $4 
		WHERE id = $5 
		RETURNING id, name, email, age, created_at, updated_at
	`

	user := &models.User{}
	err = r.db.QueryRow(
		query,
		currentUser.Name,
		currentUser.Email,
		currentUser.Age,
		currentUser.UpdatedAt,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete deletes a user
func (r *userRepository) Delete(id int) error {
	// First check if user exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, name, email, age, created_at, updated_at 
		FROM users 
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

// Count returns the total number of users
func (r *userRepository) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}