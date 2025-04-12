// internal/generator/templates/db.go - Templates for database files
package templates

// DBTemplate returns the content of the db.go file
func DBTemplate() string {
	return `// internal/db/db.go - Database connection and management
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"{{ .ModuleName }}/internal/logger"
)

// Database represents a database connection
type Database struct {
	log        logger.Logger
	connString string
	db         *sqlx.DB
}

// NewDatabase creates a new database connection
func NewDatabase(log logger.Logger, connString string) (*Database, error) {
	return &Database{
		log:        log,
		connString: connString,
	}, nil
}

// Connect connects to the database
func (d *Database) Connect() error {
	d.log.Info("Connecting to database")

	// Connect to database
	db, err := sqlx.Connect("postgres", d.connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Set database connection
	d.db = db

	d.log.Info("Connected to database")
	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		d.log.Info("Closing database connection")
		return d.db.Close()
	}
	return nil
}

// Ping pings the database
func (d *Database) Ping(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// GetDB returns the database connection
func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
`
}

// UserModelTemplate returns the template for a User model
func UserModelTemplate() string {
	return `// internal/db/models/users.go - User model
package models

import (
	"time"
)

// User represents the users table
type User struct {
	Id 	 int        ` + "`db:\"id\" json:\"id\"`" + `
	Username string ` + "`db:\"username\" json:\"username\"`" + `
	Email    string ` + "`db:\"email\" json:\"email\"`" + `
	Password string ` + "`db:\"password\" json:\"-\"`" + `
	CreatedAt time.Time  ` + "`db:\"created_at\" json:\"created_at\"`" + `
	UpdatedAt time.Time  ` + "`db:\"updated_at\" json:\"updated_at\"`" + `
}

// TableName returns the table name for User
func (u *User) TableName() string {
	return "users"
}
`
}

// DBRepositoriesTemplate returns the content of the repositories.go file
func DBRepositoriesTemplate() string {
	return `// internal/db/repositories/repositories.go - Database repositories
package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"{{ .ModuleName }}/internal/db/models"
	"{{ .ModuleName }}/internal/logger"
)

// UserRepository represents a repository for users
type UserRepository struct {
	log logger.Logger
	db  *sqlx.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(log logger.Logger, db *sqlx.DB) *UserRepository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}

// GetByID gets a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	query := "SELECT * FROM users WHERE id = $1"
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := ` + "`" + `
		INSERT INTO users (username, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	` + "`" + `

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	query := ` + "`" + `
		UPDATE users
		SET username = $1, email = $2, updated_at = $3
		WHERE id = $4
	` + "`" + `

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	now := time.Now()

	query := "UPDATE users SET WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// List lists all users
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	query := "SELECT * FROM users WHERE ORDER BY id LIMIT $1 OFFSET $2"
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}
`
}
