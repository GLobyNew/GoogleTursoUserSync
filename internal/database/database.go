package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/GLobyNew/GoogleTursoUserSync/internal/google"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Database struct {
	db *sql.DB
}

type TursoUser struct {
	PrimaryEmail string
	TgID         int64
}

func NewDatabaseConnection(databaseURL string) (*Database, error) {
	db, err := sql.Open("libsql", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &Database{db: db}, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func (db *Database) GetUserByEmail(email string) (TursoUser, error) {
	row := db.db.QueryRow("SELECT email, telegramID FROM employees WHERE email = ?", email)
	if row.Err() != nil {
		return TursoUser{}, fmt.Errorf("error querying user by email: %v", row.Err())
	}

	var user TursoUser
	err := row.Scan(&user.PrimaryEmail, &user.TgID)
	if err != nil {
		if err == sql.ErrNoRows {
			return TursoUser{}, sql.ErrNoRows
		}
		return TursoUser{}, fmt.Errorf("error scanning user row: %v", err)
	}

	return user, nil

}

func (db *Database) AddUser(user google.GoogleUser) error {
	_, err := db.db.Exec("INSERT INTO employees (email, telegramID) VALUES (?, ?)", user.PrimaryEmail, user.TgID)
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}
	return nil
}

func (db *Database) UpdateUserTgID(email string, tgID int64) error {
	_, err := db.db.Exec("UPDATE employees SET telegramID = ? WHERE email = ?", tgID, email)
	if err != nil {
		return fmt.Errorf("error updating user Telegram ID: %v", err)
	}
	return nil
}

func (db *Database) GetAllUsers() ([]TursoUser, error) {
	rows, err := db.db.Query("SELECT email, telegramID FROM employees")
	if err != nil {
		return nil, fmt.Errorf("error querying all users: %v", err)
	}
	defer rows.Close()

	var users []TursoUser
	for rows.Next() {
		var user TursoUser
		err := rows.Scan(&user.PrimaryEmail, &user.TgID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user rows: %v", err)
	}

	return users, nil
}

func userExistInTursoDB(user google.GoogleUser, tursoUsers []TursoUser) (TursoUser, bool) {
	for _, tUser := range tursoUsers {
		if user.PrimaryEmail == tUser.PrimaryEmail {
			return tUser, true
		}
	}
	return TursoUser{}, false
}

// syncUsers synchronizes Google users with the Turso database one-way,
// updating only Turso based on Google data.
func (db *Database) SyncUsers(googleUsers []google.GoogleUser, tursoUsers []TursoUser) (updated, created, upToDate int) {
	for _, gUser := range googleUsers {
		tUser, exist := userExistInTursoDB(gUser, tursoUsers)
		if exist {
			if gUser.TgID != tUser.TgID {
				// Telegram ID has changed, updating it
				log.Printf("Turos's telegram ID is not up-to-update, updating it...")
				err := db.UpdateUserTgID(gUser.PrimaryEmail, gUser.TgID)
				if err != nil {
					log.Printf("error updating TG user ID for %v in Turos: %v", gUser.PrimaryEmail, err)
					continue
				}
				updated++
				log.Printf("User %v has been successfully updated", gUser.PrimaryEmail)
			} else {
				// Telegram ID is up-to-date
				upToDate++
				log.Printf("User %v is up-to-date", gUser.PrimaryEmail)
			}
		} else {
			log.Printf("User %v is not in DB, adding him...", gUser.PrimaryEmail)
			err := db.AddUser(gUser)
			if err != nil {
				log.Printf("Error adding user %v to DB: %v", gUser.PrimaryEmail, err)
			}
			created++
			log.Printf("User %v succesfully has been added to DB", gUser.PrimaryEmail)
		}
	}
	return updated, created, upToDate
}
