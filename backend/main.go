package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	response := Response{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/submit", submitHandler)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
