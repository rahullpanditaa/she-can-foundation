package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type Response struct {
	Success bool `json:"success"`
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodPost {
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	var form FormData

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		http.Error(
			w,
			"Invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	fmt.Println("New Form Submission")
	fmt.Println("Name:", form.Name)
	fmt.Println("Email:", form.Email)
	fmt.Println("Message:", form.Message)

	// insert form data into db
	insertQuery := `
	INSERT INTO submissions (name, email, message)
	VALUES (?, ?, ?)
	`
	_, err = db.Exec(
		insertQuery,
		form.Name,
		form.Email,
		form.Message,
	)

	if err != nil {
		http.Error(
			w,
			"Database insert failed",
			http.StatusInternalServerError,
		)
		return
	}

	response := Response{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// db
var db *sql.DB

func main() {
	// initialize db
	var err error
	db, err = sql.Open("sqlite3", "./submissions.db")
	if err != nil {
		panic(err)
	}

	// create submissions table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS submissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		email TEXT,
		message TEXT
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/submit", submitHandler)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
