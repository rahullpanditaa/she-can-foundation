# She Can Foundation — Full Stack Internship Task

A full-stack web app built for the She Can Foundation internship assignment.

Users can submit contact form responses, which are stored in a SQLite database and managed through a protected admin dashboard.

---

# Features

## User Features

* Submit form with:

  * Name
  * Email
  * Message
* Form validation
* Responsive UI
* Success toast notifications
* Smooth animations

---

## Admin Features

* Admin authentication
* Protected admin dashboard
* View all submissions
* Delete submissions
* Export submissions as CSV
* Timestamped submissions
* Empty dashboard state handling

---

# Tech Stack

## Frontend

* HTML
* CSS
* JavaScript

## Backend

* Go (`net/http`)
* SQLite

---

# Project Structure

```text
she-can-foundation/
│
├── frontend/
│   ├── index.html
│   ├── style.css
│   └── script.js
│
├── backend/
│   ├── main.go
│   ├── admin.html
│   ├── login.html
│   └── submissions.db
│
└── README.md
```

---

# Setup Instructions

## 1. Clone Repository

```bash
git clone YOUR_REPOSITORY_LINK
```

---

## 2. Start Backend

Navigate into backend folder:

```bash
cd backend
```

Install SQLite driver:

```bash
go get github.com/mattn/go-sqlite3
```

Run backend server:

```bash
go run .
```

Backend runs on:

```text
http://localhost:8080
```

---

## 3. Start Frontend

Open another terminal.

Navigate into frontend folder:

```bash
cd frontend
```

Run local server:

```bash
python3 -m http.server 8000
```

Frontend runs on:

```text
http://localhost:8000
```

---

# Admin Login

Visit:

```text
http://localhost:8080/login
```

Default credentials:

```text
Username: admin1234
Password: psswrd123
```

---

# API Endpoint

## Submit Form

```http
POST /submit
```

Example request:

```json
{
  "name": "Rahul",
  "email": "rahul@example.com",
  "message": "Hello"
}
```

Example response:

```json
{
  "success": true
}
```
