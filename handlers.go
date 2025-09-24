package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Register handler
func registerHandler(c *fiber.Ctx) error {
	var req RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	if len(req.Password) < 6 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Password must be at least 6 characters long",
		})
	}

	// Create user
	user, err := createUser(&req)
	if err != nil {
		// Check for duplicate email/username
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(409).JSON(fiber.Map{
				"error": "User with this email or username already exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
		"token":   token,
	})
}

// Login handler
func loginHandler(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// Get user by email
	user, err := getUserByEmail(req.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Check password
	if !checkPasswordHash(req.Password, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   token,
	})
}

// Profile handler (protected route)
func profileHandler(c *fiber.Ctx) error {
	// Get user from context (set by auth middleware)
	userID := c.Locals("userID").(int)

	user, err := getUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// Auth middleware
func authMiddleware(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Extract token (Bearer <token>)
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid authorization format. Use: Bearer <token>",
		})
	}

	// Parse and validate token
	secret := getJWTSecret()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Store user info in context
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

// Admin middleware (requires auth middleware first)
func adminMiddleware(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	if role != "admin" && role != "super_admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}

	return c.Next()
}

// Helper function to get JWT secret
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // Default for development
	}
	return secret
}

// Admin: Create new product
func adminCreateProductHandler(c *fiber.Ctx) error {
	var req CreateProductRequest
	
	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Basic validation
	if req.Name == "" || req.Slug == "" || req.CategoryID <= 0 || req.BasePrice <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Name, slug, category_id, and base_price are required",
		})
	}

	// Create product
	product, err := createProduct(&req)
	if err != nil {
		// Check for duplicate slug
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(409).JSON(fiber.Map{
				"error": "Product with this slug already exists",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create product",
			"details": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Product created successfully",
		"product": product,
	})
}

// Admin: Update existing product
func adminUpdateProductHandler(c *fiber.Ctx) error {
	// Get product ID from URL
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req UpdateProductRequest
	
	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update product
	product, err := updateProduct(id, &req)
	if err != nil {
		if err.Error() == "no fields to update" {
			return c.Status(400).JSON(fiber.Map{
				"error": "No fields provided for update",
			})
		}
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(409).JSON(fiber.Map{
				"error": "Product with this slug already exists",
			})
		}
		if err.Error() == "sql: no rows in result set" {
			return c.Status(404).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update product",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product updated successfully",
		"product": product,
	})
}

// Admin: Delete product (soft delete)
func adminDeleteProductHandler(c *fiber.Ctx) error {
	// Get product ID from URL
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Delete product (soft delete)
	err = deleteProduct(id)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(404).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete product",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// Admin: Get all products (including inactive)
func adminGetProductsHandler(c *fiber.Ctx) error {
	query := `
		SELECT id, name, slug, description, category_id, base_price, is_active, is_featured 
		FROM products 
		ORDER BY created_at DESC
	`
	
	rows, err := db.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch products",
			"details": err.Error(),
		})
	}
	defer rows.Close()
	
	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID, &p.BasePrice, &p.IsActive, &p.IsFeatured)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to scan products",
				"details": err.Error(),
			})
		}
		products = append(products, p)
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    len(products),
		"message": "All products (including inactive)",
	})
}
