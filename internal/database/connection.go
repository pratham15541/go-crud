package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pratham15541/go-crud/internal/config"
	_ "github.com/lib/pq"
)

// NewConnection creates a new database connection
func NewConnection(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

// RunMigrations runs the database migrations
func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		age INTEGER CHECK (age > 0 AND age < 150),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create updated_at trigger function
	updatedAtTrigger := `
	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = CURRENT_TIMESTAMP;
		RETURN NEW;
	END;
	$$ language 'plpgsql';
	`

	if _, err := db.Exec(updatedAtTrigger); err != nil {
		return fmt.Errorf("failed to create updated_at trigger function: %w", err)
	}

	// Create trigger for users table
	usersTrigger := `
	DROP TRIGGER IF EXISTS update_users_updated_at ON users;
	CREATE TRIGGER update_users_updated_at
		BEFORE UPDATE ON users
		FOR EACH ROW
		EXECUTE FUNCTION update_updated_at_column();
	`

	if _, err := db.Exec(usersTrigger); err != nil {
		return fmt.Errorf("failed to create users trigger: %w", err)
	}

	// Create index on email for faster lookups
	emailIndex := `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`
	if _, err := db.Exec(emailIndex); err != nil {
		return fmt.Errorf("failed to create email index: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}