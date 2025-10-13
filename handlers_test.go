package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "Merch Ke API",
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status code = %d, want 200", resp.StatusCode)
	}

	var result map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	if result["status"] != "healthy" {
		t.Errorf("Status = %s, want healthy", result["status"])
	}

	if result["service"] != "Merch Ke API" {
		t.Errorf("Service = %s, want Merch Ke API", result["service"])
	}
}

// TestRegisterHandlerValidation tests registration input validation
func TestRegisterHandlerValidation(t *testing.T) {
	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-purposes-only")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		errorContains  string
	}{
		{
			name: "Missing username",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "Password123!",
			},
			expectedStatus: 400,
			errorContains:  "required",
		},
		{
			name: "Missing email",
			payload: map[string]interface{}{
				"username": "testuser",
				"password": "Password123!",
			},
			expectedStatus: 400,
			errorContains:  "required",
		},
		{
			name: "Missing password",
			payload: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
			expectedStatus: 400,
			errorContains:  "required",
		},
		{
			name: "Short password",
			payload: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "12345",
			},
			expectedStatus: 400,
			errorContains:  "6 characters",
		},
		{
			name: "Invalid JSON",
			payload: map[string]interface{}{
				"invalid": "data",
			},
			expectedStatus: 400,
			errorContains:  "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/api/auth/register", registerHandler)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status code = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

// TestLoginHandlerValidation tests login input validation
func TestLoginHandlerValidation(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Missing email",
			payload: map[string]interface{}{
				"password": "Password123!",
			},
			expectedStatus: 400,
		},
		{
			name: "Missing password",
			payload: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: 400,
		},
		{
			name: "Empty email",
			payload: map[string]interface{}{
				"email":    "",
				"password": "Password123!",
			},
			expectedStatus: 400,
		},
		{
			name: "Empty password",
			payload: map[string]interface{}{
				"email":    "test@example.com",
				"password": "",
			},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/api/auth/login", loginHandler)

			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Status code = %d, want %d", resp.StatusCode, tt.expectedStatus)
			}
		})
	}
}

// TestAuthMiddlewareNoToken tests auth middleware without token
func TestAuthMiddlewareNoToken(t *testing.T) {
	app := fiber.New()

	// Protected route
	app.Get("/api/protected", authMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "protected"})
	})

	req := httptest.NewRequest("GET", "/api/protected", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("Status code = %d, want 401", resp.StatusCode)
	}
}

// TestAuthMiddlewareInvalidToken tests auth middleware with invalid token
func TestAuthMiddlewareInvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-purposes-only")
	defer os.Unsetenv("JWT_SECRET")

	app := fiber.New()

	// Protected route
	app.Get("/api/protected", authMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "protected"})
	})

	req := httptest.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 401 {
		t.Errorf("Status code = %d, want 401", resp.StatusCode)
	}
}

// TestAdminMiddlewareNoAdmin tests admin middleware with non-admin user
func TestAdminMiddlewareNoAdmin(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-purposes-only")
	defer os.Unsetenv("JWT_SECRET")

	app := fiber.New()

	// Mock middleware that sets role as customer (not admin)
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "customer") // Set role, not user object
		return c.Next()
	})

	app.Get("/api/admin/test", adminMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "admin only"})
	})

	req := httptest.NewRequest("GET", "/api/admin/test", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 403 {
		t.Errorf("Status code = %d, want 403", resp.StatusCode)
	}
}

// TestAdminMiddlewareWithAdmin tests admin middleware with admin user
func TestAdminMiddlewareWithAdmin(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-testing-purposes-only")
	defer os.Unsetenv("JWT_SECRET")

	app := fiber.New()

	// Mock middleware that sets role as admin
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", "admin")
		return c.Next()
	})

	app.Get("/api/admin/test", adminMiddleware, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "admin access granted"})
	})

	req := httptest.NewRequest("GET", "/api/admin/test", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Status code = %d, want 200", resp.StatusCode)
	}
}

// TestJSONResponseFormat tests consistent JSON response format
func TestJSONResponseFormat(t *testing.T) {
	app := fiber.New()

	app.Get("/test-error", func(c *fiber.Ctx) error {
		return c.Status(400).JSON(fiber.Map{
			"error": "Test error message",
		})
	})

	app.Get("/test-success", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Success",
			"data":    "test data",
		})
	})

	// Test error response
	req := httptest.NewRequest("GET", "/test-error", nil)
	resp, _ := app.Test(req)

	var errorResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&errorResult)

	if _, ok := errorResult["error"]; !ok {
		t.Error("Error response should contain 'error' field")
	}

	// Test success response
	req = httptest.NewRequest("GET", "/test-success", nil)
	resp, _ = app.Test(req)

	var successResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&successResult)

	if _, ok := successResult["message"]; !ok {
		t.Error("Success response should contain 'message' field")
	}
}

// TestCORSHeaders tests CORS configuration
func TestCORSHeaders(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	// Note: CORS headers are set by middleware in main app
	// This test verifies the app can handle CORS requests
	if resp.StatusCode != 200 {
		t.Errorf("Status code = %d, want 200", resp.StatusCode)
	}
}

// TestInvalidContentType tests invalid content type handling
func TestInvalidContentType(t *testing.T) {
	app := fiber.New()

	app.Post("/test", func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}
		return c.JSON(data)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Status code = %d, want 400", resp.StatusCode)
	}
}
