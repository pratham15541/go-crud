package repository

import (
	"github.com/pratham15541/go-crud/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.CreateUserRequest) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetAll(limit, offset int) ([]*models.User, error)
	Update(id int, user *models.UpdateUserRequest) (*models.User, error)
	Delete(id int) error
	GetByEmail(email string) (*models.User, error)
	Count() (int64, error)
}

// HealthRepository defines the interface for health check operations
type HealthRepository interface {
	Ping() error
}