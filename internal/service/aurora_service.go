package service

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"os"
)

type AuroraService struct {
	db *sql.DB
}

// NewAuroraService initializes the MySQL connection.
func NewAuroraService() (*AuroraService, error) {

	// Read connection details from environment variables
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")

	// Build the connection string
	connStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, host, database)
	fmt.Println("Connecting to MySQL with:", connStr)

	// Open the database connection
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %v", err)
	}

	return &AuroraService{db: db}, nil
}

// GetValueByKey retrieves a value for a given key from the MySQL database.
func (s *AuroraService) GetValueByKey(key string) (string, error) {
	var value string
	query := "SELECT value_field FROM key_value_store WHERE key_field = ?"

	err := s.db.QueryRow(query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("key not found")
		}
		return "", fmt.Errorf("query failed: %v", err)
	}

	return value, nil
}
