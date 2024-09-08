package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
	"golang.org/x/crypto/bcrypt"    // Импортируем пакет для хеширования паролей
)

// InitDatabase инициализирует базу данных и создает таблицы
func InitDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "bot_data.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

// RegisterUser регистрирует нового пользователя
func RegisterUser(username, password string, db *sql.DB) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if user exists: %v", err)
		return false
	}
	if exists {
		log.Printf("User %s already exists.", username)
		return false
	}

	// Хешируем пароль перед сохранением
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return false
	}

	// Вставляем нового пользователя с хешированным паролем
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return false
	}
	return true
}

// IsValidUser проверяет, действителен ли пользователь
func IsValidUser(username, password string, db *sql.DB) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error checking user: %v", err)
		return false
	}

	// Сравниваем введенный пароль с хешированным паролем
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Invalid password for user %s: %v", username, err)
		return false
	}
	return true
}
