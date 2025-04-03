package db

import (
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres_password")
	os.Setenv("DB_NAME", "saas_db")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")

	db, err := NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}
