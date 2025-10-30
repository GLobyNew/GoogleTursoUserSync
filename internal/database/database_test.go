package database

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByEmail(t *testing.T) {
	godotenv.Load()
	databaseURL := os.Getenv("DATABASE_TEST_URL")
	if databaseURL == "" {
		t.Fatalf("DATABASE_TEST_URL is not set")
	}

	db, err := NewDatabaseConnection(databaseURL)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	testEmail := os.Getenv("DATABASE_TEST_USER_EMAIL")
	if testEmail == "" {
		t.Fatalf("DATABASE_TEST_USER_EMAIL is not set")
	}

	user, err := db.GetUserByEmail(testEmail)
	if err != nil {
		t.Fatalf("Error getting user by email: %v", err)
	}

	t.Logf("Retrieved user: %+v", user)
}

func TestGetUserByEmail_NoUser(t *testing.T) {
	godotenv.Load()
	databaseURL := os.Getenv("DATABASE_TEST_URL")
	if databaseURL == "" {
		t.Fatalf("DATABASE_TEST_URL is not set")
	}

	db, err := NewDatabaseConnection(databaseURL)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	nonExistentEmail := "test@test.com"
	_, err = db.GetUserByEmail(nonExistentEmail)
	assert.EqualValues(t, err, sql.ErrNoRows)
}
