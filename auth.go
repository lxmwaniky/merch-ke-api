package main

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User struct for authentication
type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"` // Never return password in JSON
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Phone         string    `json:"phone"`
	Role          string    `json:"role"`
	IsActive      bool      `json:"is_active"`
	EmailVerified bool      `json:"email_verified"`
	WalletBalance float64   `json:"wallet_balance"`
	CreatedAt     time.Time `json:"created_at"`
}

// WalletTransaction for tracking token movements
type WalletTransaction struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	OrderID      *int      `json:"order_id,omitempty"`
	Amount       float64   `json:"amount"`
	Type         string    `json:"type"` // "credit" or "debit"
	Description  string    `json:"description"`
	BalanceAfter float64   `json:"balance_after"`
	CreatedAt    time.Time `json:"created_at"`
}

// Registration request struct
type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

// Login request struct
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// JWT Claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Hash password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Check if password matches hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate JWT token
func generateJWT(user *User) (string, error) {
	// Get JWT secret from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // Default for development
	}

	// Create claims
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "merch-ke-api",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Create new user in database
func createUser(req *RegisterRequest) (*User, error) {
	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Insert user into database
	query := `
		INSERT INTO auth.users (username, email, password_hash, first_name, last_name, phone, role)
		VALUES ($1, $2, $3, $4, $5, $6, 'customer')
		RETURNING id, username, email, first_name, last_name, phone, role, is_active, email_verified, created_at
	`

	var user User
	err = db.QueryRow(query, req.Username, req.Email, hashedPassword, req.FirstName, req.LastName, req.Phone).
		Scan(&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.Phone, &user.Role, &user.IsActive, &user.EmailVerified, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by email for login
func getUserByEmail(email string) (*User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, phone, role, is_active, email_verified, wallet_balance, created_at
		FROM auth.users 
		WHERE email = $1 AND is_active = true
	`

	var user User
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.Role,
		&user.IsActive, &user.EmailVerified, &user.WalletBalance, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Get user by ID
func getUserByID(userID int) (*User, error) {
	query := `
		SELECT id, username, email, first_name, last_name, phone, role, is_active, email_verified, wallet_balance, created_at
		FROM auth.users 
		WHERE id = $1 AND is_active = true
	`

	var user User
	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.Phone, &user.Role,
		&user.IsActive, &user.EmailVerified, &user.WalletBalance, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Get wallet balance
func getWalletBalance(userID int) (float64, error) {
	var balance float64
	query := `SELECT wallet_balance FROM auth.users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(&balance)
	return balance, err
}

// Add wallet transaction and update balance
func addWalletTransaction(userID int, amount float64, transactionType, description string, orderID *int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get current balance
	var currentBalance float64
	err = tx.QueryRow(`SELECT wallet_balance FROM auth.users WHERE id = $1 FOR UPDATE`, userID).Scan(&currentBalance)
	if err != nil {
		return err
	}

	// Calculate new balance
	newBalance := currentBalance
	if transactionType == "credit" {
		newBalance += amount
	} else {
		newBalance -= amount
		if newBalance < 0 {
			return errors.New("insufficient wallet balance")
		}
	}

	// Update user balance
	_, err = tx.Exec(`UPDATE auth.users SET wallet_balance = $1 WHERE id = $2`, newBalance, userID)
	if err != nil {
		return err
	}

	// Insert transaction record
	_, err = tx.Exec(`
		INSERT INTO auth.wallet_transactions (user_id, order_id, amount, type, description, balance_after)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, orderID, amount, transactionType, description, newBalance)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Get wallet transactions
func getWalletTransactions(userID int, limit int) ([]WalletTransaction, error) {
	query := `
		SELECT id, user_id, order_id, amount, type, description, balance_after, created_at
		FROM auth.wallet_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []WalletTransaction
	for rows.Next() {
		var t WalletTransaction
		err := rows.Scan(&t.ID, &t.UserID, &t.OrderID, &t.Amount, &t.Type, &t.Description, &t.BalanceAfter, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
