package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	DB_CONNECTION = "postgres://postgres:xpertscan@localhost:5432/arunhiremath?sslmode=disable"
)

type User struct {
	Age  int
	Name string
}

type ConnectionHandler struct {
	db *sql.DB
}

func ModelsInstance() (*ConnectionHandler, error) {
	db, err := sql.Open("postgres", DB_CONNECTION)
	if err != nil {
		fmt.Println("failed to open the db connection", err.Error())
		return nil, fmt.Errorf("failed to connect database %s", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		fmt.Println("failed to ping the db connection")
		return nil, fmt.Errorf("failed to connect to the database %s", err)
	}
	fmt.Println("connected to database successfully!")

	return &ConnectionHandler{
		db: db,
	}, nil
}

func (conn *ConnectionHandler) GetUserDetails() []User {
	rows, err := conn.db.Query("select * from users")
	userData := []User{}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println("\nUsers:")
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Age, &user.Name)
		if err != nil {
			log.Fatal(err)
		}
		userData = append(userData, user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return userData
}

func (conn *ConnectionHandler) AddUserDetails(age int, name string) error {
	sqlStatement := `
		INSERT INTO users (age, name)
		VALUES ($1, $2)` // PostgreSQL specific clause to return the ID

	// QueryRow is used when you expect a single row back (like an ID)
	conn.db.QueryRow(sqlStatement, age, name)
	return nil

}
