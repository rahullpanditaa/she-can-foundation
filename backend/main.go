package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// hardcoded for now
const ADMIN_USERNAME = "admin1234"
const ADMIN_PASSWORD = "psswrd123"

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type Response struct {
	Success bool `json:"success"`
}

type Submission struct {
	ID        int
	Name      string
	Email     string
	Message   string
	CreatedAt string
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

func adminHandler(w http.ResponseWriter, r *http.Request) {

	// use auth middleware to protect admin dashboard
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := db.Query(`
		SELECT id, name, email, message, created_at
		FROM submissions
		ORDER BY id DESC
	`)

	if err != nil {
		http.Error(
			w,
			"Failed to fetch submissions",
			http.StatusInternalServerError,
		)
		return
	}
	defer rows.Close()

	var submissions []Submission
	for rows.Next() {
		var submission Submission

		err := rows.Scan(
			&submission.ID,
			&submission.Name,
			&submission.Email,
			&submission.Message,
			&submission.CreatedAt,
		)

		if err != nil {
			continue
		}
		submissions = append(submissions, submission)
	}

	tmpl := template.Must(
		template.ParseFiles("admin.html"),
	)
	tmpl.Execute(w, submissions)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {

	// protect route
	if !isAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	id := r.URL.Query().Get("id")
	_, err := db.Exec(
		"DELETE FROM submissions WHERE id = ?",
		id,
	)

	if err != nil {
		http.Error(
			w,
			"Delete failed",
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		w,
		r,
		"/admin",
		http.StatusSeeOther,
	)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(
			template.ParseFiles("login.html"),
		)
		tmpl.Execute(w, nil)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == ADMIN_USERNAME &&
		password == ADMIN_PASSWORD {
		cookie := http.Cookie{
			Name:  "authenticated",
			Value: "true",
			Path:  "/",
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}
	tmpl := template.Must(template.ParseFiles("login.html"))
	tmpl.Execute(w, "Invalid username or password")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "authenticated",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// auth middleware
func isAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("authenticated")
	if err != nil {
		return false
	}
	return cookie.Value == "true"
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
		message TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
