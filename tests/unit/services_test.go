package unit

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/pratham15541/go-crud/internal/models"
	"github.com/pratham15541/go-crud/internal/services"
)

// MockUserRepository implements UserRepository interface for testing
type MockUserRepository struct {
	users map[int]*models.User
	nextID int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

func (m *MockUserRepository) Create(req *models.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		ID:    m.nextID,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}
	m.users[m.nextID] = user
	m.nextID++
	return user, nil
}

func (m *MockUserRepository) GetByID(id int) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) GetAll(limit, offset int) ([]*models.User, error) {
	var users []*models.User
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) Update(id int, req *models.UpdateUserRequest) (*models.User, error) {
	if user, exists := m.users[id]; exists {
		if req.Name != "" {
			user.Name = req.Name
		}
		if req.Email != "" {
			user.Email = req.Email
		}
		if req.Age != 0 {
			user.Age = req.Age
		}
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Delete(id int) error {
	if _, exists := m.users[id]; exists {
		delete(m.users, id)
		return nil
	}
	return fmt.Errorf("user not found")
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Count() (int64, error) {
	return int64(len(m.users)), nil
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	req := &models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	user, err := userService.CreateUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, 30, user.Age)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	req := &models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	// Create first user
	_, err := userService.CreateUser(req)
	assert.NoError(t, err)

	// Try to create user with same email
	_, err = userService.CreateUser(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	// Create a user first
	req := &models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	createdUser, _ := userService.CreateUser(req)

	// Get the user
	user, err := userService.GetUser(createdUser.ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, createdUser.ID, user.ID)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	userService := services.NewUserService(mockRepo)

	_, err := userService.GetUser(999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}