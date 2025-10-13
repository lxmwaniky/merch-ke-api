package main

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword tests password hashing
func TestHashPassword(t *testing.T) {
	password := "TestPassword123!"

	hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	if hash == "" {
		t.Error("Hash should not be empty")
	}

	if hash == password {
		t.Error("Hash should not equal plain text password")
	}

	// Verify we can check the password
	if !checkPasswordHash(password, hash) {
		t.Error("Password hash verification failed")
	}
}

// TestCheckPasswordHash tests password verification
func TestCheckPasswordHash(t *testing.T) {
	password := "CorrectPassword123!"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     string(hash),
			want:     true,
		},
		{
			name:     "Wrong password",
			password: "WrongPassword",
			hash:     string(hash),
			want:     false,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     string(hash),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkPasswordHash(tt.password, tt.hash)
			if result != tt.want {
				t.Errorf("checkPasswordHash() = %v, want %v", result, tt.want)
			}
		})
	}
}

// TestGenerateJWT tests JWT token generation
func TestGenerateJWT(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-purposes-only")
	defer os.Unsetenv("JWT_SECRET")

	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "customer",
	}

	token, err := generateJWT(user)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	// Parse and verify token
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Token should be valid")
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		t.Fatal("Failed to extract claims from token")
	}

	if claims.UserID != user.ID {
		t.Errorf("UserID = %d, want %d", claims.UserID, user.ID)
	}

	if claims.Username != user.Username {
		t.Errorf("Username = %s, want %s", claims.Username, user.Username)
	}

	if claims.Email != user.Email {
		t.Errorf("Email = %s, want %s", claims.Email, user.Email)
	}

	if claims.Role != user.Role {
		t.Errorf("Role = %s, want %s", claims.Role, user.Role)
	}
}

// TestGenerateJWTNoSecret tests JWT generation without secret
func TestGenerateJWTNoSecret(t *testing.T) {
	// Ensure JWT_SECRET is not set
	os.Unsetenv("JWT_SECRET")

	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "customer",
	}

	_, err := generateJWT(user)
	if err == nil {
		t.Error("Expected error when JWT_SECRET is not set")
	}
}

// TestParseJWT tests JWT token parsing
func TestParseJWT(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-purposes-only")
	defer os.Unsetenv("JWT_SECRET")

	// Create a valid token
	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "customer",
	}

	validToken, _ := generateJWT(user)

	// Parse the token
	parsedToken, err := jwt.ParseWithClaims(validToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse valid token: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Valid token should be marked as valid")
	}

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	if claims.UserID != user.ID {
		t.Errorf("UserID = %d, want %d", claims.UserID, user.ID)
	}
}

// TestUserStruct tests the User struct
func TestUserStruct(t *testing.T) {
	user := User{
		ID:            1,
		Username:      "testuser",
		Email:         "test@example.com",
		PasswordHash:  "hashed_password",
		FirstName:     "Test",
		LastName:      "User",
		Phone:         "+254712345678",
		Role:          "customer",
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     time.Now(),
	}

	if user.ID != 1 {
		t.Errorf("ID = %d, want 1", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("Username = %s, want testuser", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Email = %s, want test@example.com", user.Email)
	}

	if user.Role != "customer" {
		t.Errorf("Role = %s, want customer", user.Role)
	}

	if !user.IsActive {
		t.Error("IsActive should be true")
	}

	if user.EmailVerified {
		t.Error("EmailVerified should be false")
	}
}
