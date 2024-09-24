package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	log.Println("Connecting to database")

	db, err = sql.Open("postgres", "user=postgres dbname=test password=mysecretpassword sslmode=disable")

	if err != nil {
		log.Fatalf("Connection to database failed with error: %s", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	} else {
		createTableIfNotExists()
		log.Println("Database connection successful")
	}

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/create", createUser)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Create the table if it doesn't exist
func createTableIfNotExists() {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100)
    );`

	_, err := db.Exec(query)
	if err != nil {
		log.Printf("failed to create users table: %v", err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT name FROM users")
	if err != nil {
		log.Printf("Error retreiving users: %v", err)
		fmt.Fprintf(w, "Error retrieving users")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue // Skip to next row on error
		}
		fmt.Fprintf(w, "User: %s\n", name)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("name")
	if username == "" {
		fmt.Fprintf(w, "Missing name in parameters")
		return
	}

	stmt, err := db.Prepare("INSERT INTO users (name) VALUES ($1)")
	if err != nil {
		log.Printf("Error preparing create user statement: %v", err)
		fmt.Fprintf(w, "Error creating user")
		return
	}
	defer stmt.Close() // Close prepared statement after use

	_, err = stmt.Exec(username)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		fmt.Fprintf(w, "Error creating user")
		return
	}

	fmt.Fprintf(w, "User %s created successfully", username)
}
