package integration

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pratham15541/go-crud/internal/config"
	"github.com/pratham15541/go-crud/internal/database"
	"github.com/pratham15541/go-crud/internal/handlers"
	"github.com/pratham15541/go-crud/internal/models"
	"github.com/pratham15541/go-crud/internal/repository"
	"github.com/pratham15541/go-crud/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	_ "github.com/lib/pq"
)

type IntegrationTestSuite struct {
	suite.Suite
	db     *sql.DB
	router *mux.Router
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Load test environment
	godotenv.Load("../../.env.example")
	os.Setenv("DB_NAME", "crud_demo_test")

	// Load configuration
	cfg := config.Load()

	// Connect to postgres to create test database
	testDbConfig := cfg.Database
	testDbConfig.Name = "postgres"
	pgDb, err := database.NewConnection(testDbConfig)
	suite.Require().NoError(err)

	// Create test database
	_, err = pgDb.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", cfg.Database.Name))
	suite.Require().NoError(err)
	_, err = pgDb.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.Database.Name))
	suite.Require().NoError(err)
	pgDb.Close()

	// Connect to test database
	db, err := database.NewConnection(cfg.Database)
	suite.Require().NoError(err)
	suite.db = db

	// Run migrations
	err = database.RunMigrations(db)
	suite.Require().NoError(err)

	// Setup router
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	healthHandler := handlers.NewHealthHandler(db)

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	// Health check
	api.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")

	// User routes
	userRoutes := api.PathPrefix("/users").Subrouter()
	userRoutes.HandleFunc("", userHandler.GetUsers).Methods("GET")
	userRoutes.HandleFunc("", userHandler.CreateUser).Methods("POST")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.GetUser).Methods("GET")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")

	suite.router = router
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Clean up users table before each test
	_, err := suite.db.Exec("DELETE FROM users")
	suite.Require().NoError(err)
}

func (suite *IntegrationTestSuite) TestHealthCheck() {
	req, _ := http.NewRequest("GET", "/api/v1/health", nil)
	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var response models.HealthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "healthy", response.Status)
}

func (suite *IntegrationTestSuite) TestCreateUser() {
	user := models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}

	jsonUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusCreated, rr.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "User created successfully", response.Message)
}

func (suite *IntegrationTestSuite) TestGetUsers() {
	// Create a test user first
	user := models.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	jsonUser, _ := json.Marshal(user)
	createReq, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(jsonUser))
	createReq.Header.Set("Content-Type", "application/json")
	createRr := httptest.NewRecorder()
	suite.router.ServeHTTP(createRr, createReq)

	// Get users
	getReq, _ := http.NewRequest("GET", "/api/v1/users", nil)
	getRr := httptest.NewRecorder()

	suite.router.ServeHTTP(getRr, getReq)

	assert.Equal(suite.T(), http.StatusOK, getRr.Code)

	var response models.SuccessResponse
	err := json.Unmarshal(getRr.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}