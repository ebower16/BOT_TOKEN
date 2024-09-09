package database

import (
	"database/sql"
	"fmt"
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
        password TEXT NOT NULL,
        chat_id INTEGER,
        registration_time DATETIME DEFAULT CURRENT_TIMESTAMP,
        message TEXT,
        money REAL DEFAULT 0
    );`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

// RegisterUser регистрирует нового пользователя
func RegisterUser(username, password string, db *sql.DB, chatID int64) bool {
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

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return false
	}

	_, err = db.Exec("INSERT INTO users (username, password, chat_id) VALUES (?, ?, ?)", username, hashedPassword, chatID)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return false
	}

	log.Printf("User %s registered successfully.", username)
	return true
}

// CheckUserCredentials проверяет логин и пароль пользователя
func CheckUserCredentials(username, password string, db *sql.DB) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Printf("Invalid password for user %s: %v", username, err)
		return false
	}

	log.Printf("User %s logged in successfully.", username)
	return true
}

// GetAllUsers возвращает список всех пользователей
func GetAllUsers(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT username, chat_id, registration_time, message, money FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		var chatID int64
		var registrationTime string
		var message string
		var money float64

		if err := rows.Scan(&username, &chatID, &registrationTime, &message, &money); err != nil {
			return nil, err
		}
		userInfo := fmt.Sprintf("Username: %s, Chat ID: %d, Registration Time: %s, Message: %s, Money: %.2f", username, chatID, registrationTime, message, money)
		users = append(users, userInfo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
