package cmd

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var DB = setupDB()

func setupDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./godirectory.db")
	checkErr(err)

	return db
}

func Init() {
	var err error

	username := "john"
	email := "john@example.com"
	password := "supersecret"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	checkErr(err)

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
                 id INT,
                 user_name VARCHAR(60),  
                 user_email VARCHAR(60),  
                 user_password VARCHAR(60),  
                 user_created TIMESTAMP WITH TIME ZONE,
                 user_last_login TIMESTAMP WITH TIME ZONE, 
                 PRIMARY KEY  (id),  
                 CONSTRAINT users_email UNIQUE (user_email)
            );`)

	checkErr(err)

	err = DB.QueryRow("SELECT user_email FROM users WHERE user_name = $1", username).Scan(&email)
	if err == sql.ErrNoRows {
		_, err = DB.Exec(`INSERT INTO users (user_name, user_email, user_password) VALUES ($1, $2, $3);`, username, email, hashedPassword)
	} else if err != nil {
		log.Print(err)
	}

	checkErr(err)
}
