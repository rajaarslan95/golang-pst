package db

import (
	"user-manager/helper"
	"user-manager/schemas"

	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	addr string
	conn *sql.DB
}

func NewDBManager() *PostgresDB {
	host := helper.Getenv("DATABASE_HOST", "localhost")
	username := helper.Getenv("DATABASE_USERNAME", "postgres")
	password := helper.Getenv("DATABASE_PASSWORD", "mysecretpassword")
	port := helper.Getenv("DATABASE_PORT", "5432")
	database := helper.Getenv("DATABASE_DBNAME", "")

	addr := fmt.Sprintf("sslmode=disable user=%s password=%s host=%s port=%s dbname=%s", username, password, host, port, database)

	return &PostgresDB{
		addr: addr,
	}
}

func (m *PostgresDB) Connect() {
	var err error

	log.Printf("[PostgresDB]: Connecting to database at %s...", m.addr)

	m.conn, err = sql.Open("postgres", m.addr)

	if err != nil {
		log.Fatalf("[PostgresDB] Connection failure: %s", err)
	}

	err = m.conn.Ping()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	} else {
		createTableIfNotExists(m.conn)
		log.Println("Database connection successful")
	}
}

// Create the table if it doesn't exist
func createTableIfNotExists(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100),
        age INT
    );`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("failed to create users table: %v", err)
	}
}

func (s *PostgresDB) AddUser(user schemas.User) error {
	_, err := s.conn.Exec("INSERT INTO users (name, email, age) VALUES ($1, $2, $3)", user.Name, user.Email, user.Age)
	return err
}

// Get User handler
func (s *PostgresDB) GetUser(id int) (schemas.User, error) {
	var user schemas.User
	err := s.conn.QueryRow("SELECT id, name, email, age FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Email, &user.Age)

	if err != nil {
		log.Printf("[PostgresDB]: Error scanning row: %s", err)
	}
	return user, err
}

// Update User handler
func (s *PostgresDB) UpdateUser(user schemas.User) error {
	_, err := s.conn.Exec("UPDATE users SET name=$1, email=$2, age=$3 WHERE id=$4", user.Name, user.Email, user.Age, user.ID)
	if err != nil {
		log.Printf("[PostgresDB]: Error updating row: %s", err)
	}

	return err
}

// Delete User handler
func (s *PostgresDB) DeleteUser(id int) error {
	_, err := s.conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		log.Printf("[PostgresDB]: Error deleting row: %s", err)
	}
	return err
}
