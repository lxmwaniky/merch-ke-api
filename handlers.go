package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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

// Optional auth middleware (allows both authenticated and guest users)
func optionalAuthMiddleware(c *fiber.Ctx) error {
	// Get token from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// No token provided - proceed as guest user
		return c.Next()
	}

	// Extract token (Bearer <token>)
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
	if tokenString == authHeader {
		// Invalid format - proceed as guest user
		return c.Next()
	}

	// Parse and validate token
	secret := getJWTSecret()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		// Invalid token - proceed as guest user
		return c.Next()
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok {
		// Invalid claims - proceed as guest user
		return c.Next()
	}

	// Store user info in context (user is authenticated)
	c.Locals("user", claims)
	c.Locals("userID", claims.UserID)
	c.Locals("username", claims.Username)
	c.Locals("email", claims.Email)
	c.Locals("role", claims.Role)

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
			"error":   "Failed to create product",
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
			"error":   "Failed to update product",
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
			"error":   "Failed to delete product",
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
			"error":   "Failed to fetch products",
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
				"error":   "Failed to scan products",
				"details": err.Error(),
			})
		}
		products = append(products, p)
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    len(products),
		"message":  "All products (including inactive)",
	})
}

// =====================================================
// CART HANDLERS
// =====================================================

// Add item to cart (works for both authenticated and guest users)
func addToCartHandler(c *fiber.Ctx) error {
	var req AddToCartRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.ProductID <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Valid product_id is required",
		})
	}

	if req.Quantity <= 0 {
		req.Quantity = 1 // Default to 1
	}

	// Check if user is authenticated
	user := c.Locals("user")

	if user != nil {
		// Authenticated user - use user cart
		userClaims := user.(*Claims)
		userID := userClaims.UserID

		err := addToUserCart(userID, req.ProductID, req.Quantity)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to add item to cart",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Item added to cart successfully",
		})
	} else {
		// Guest user - use session cart
		sessionID := c.Get("X-Session-ID", "")
		if sessionID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Session ID required for guest users (send X-Session-ID header)",
			})
		}

		err := addToGuestCart(sessionID, req.ProductID, req.Quantity)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to add item to cart",
				"details": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Item added to cart successfully",
		})
	}
}

// Get cart items and summary
func getCartHandler(c *fiber.Ctx) error {
	user := c.Locals("user")

	if user != nil {
		// Authenticated user
		userClaims := user.(*Claims)
		userID := userClaims.UserID

		summary, err := getCartSummary(&userID, nil)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to get cart",
				"details": err.Error(),
			})
		}

		return c.JSON(summary)
	} else {
		// Guest user
		sessionID := c.Get("X-Session-ID", "")
		if sessionID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Session ID required for guest users (send X-Session-ID header)",
			})
		}

		summary, err := getCartSummary(nil, &sessionID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to get cart",
				"details": err.Error(),
			})
		}

		return c.JSON(summary)
	}
}

// Update cart item quantity
func updateCartHandler(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	var req struct {
		Quantity int `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user := c.Locals("user")

	if user != nil {
		// Authenticated user
		userClaims := user.(*Claims)
		userID := userClaims.UserID

		err := updateCartItemQuantity(&userID, nil, productID, req.Quantity)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to update cart item",
				"details": err.Error(),
			})
		}
	} else {
		// Guest user
		sessionID := c.Get("X-Session-ID", "")
		if sessionID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Session ID required for guest users",
			})
		}

		err := updateCartItemQuantity(nil, &sessionID, productID, req.Quantity)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to update cart item",
				"details": err.Error(),
			})
		}
	}

	message := "Cart item updated successfully"
	if req.Quantity <= 0 {
		message = "Item removed from cart"
	}

	return c.JSON(fiber.Map{
		"message": message,
	})
}

// Remove item from cart
func removeFromCartHandler(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("productId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	user := c.Locals("user")

	if user != nil {
		// Authenticated user
		userClaims := user.(*Claims)
		userID := userClaims.UserID

		err := removeFromCart(&userID, nil, productID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to remove item from cart",
				"details": err.Error(),
			})
		}
	} else {
		// Guest user
		sessionID := c.Get("X-Session-ID", "")
		if sessionID == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Session ID required for guest users",
			})
		}

		err := removeFromCart(nil, &sessionID, productID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error":   "Failed to remove item from cart",
				"details": err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "Item removed from cart successfully",
	})
}

// Get user points (authenticated users only)
func getUserPointsHandler(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	userClaims := user.(*Claims)
	userID := userClaims.UserID

	points, err := getUserPoints(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get user points",
			"details": err.Error(),
		})
	}

	return c.JSON(points)
}

// Migrate guest cart to user cart (called during login/register)
func migrateCartHandler(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	sessionID := c.Get("X-Session-ID", "")
	if sessionID == "" {
		return c.JSON(fiber.Map{
			"message": "No guest cart to migrate",
		})
	}

	userClaims := user.(*Claims)
	userID := userClaims.UserID

	err := migrateGuestCartToUser(sessionID, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to migrate cart",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Guest cart migrated successfully",
	})
}

// =====================================================
// ORDER HANDLERS
// =====================================================

// Create order from cart
func createOrderHandler(c *fiber.Ctx) error {
	var req CreateOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user info (authenticated or guest)
	user := c.Locals("user")
	sessionID := c.Get("X-Session-ID", "")

	var userID *int
	var sessionIDPtr *string

	if user != nil {
		// Authenticated user
		userClaims := user.(*Claims)
		userID = &userClaims.UserID
	} else if sessionID != "" {
		// Guest user
		sessionIDPtr = &sessionID
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "Either authentication or session ID required",
		})
	}

	// Create order
	order, err := createOrderFromCart(userID, sessionIDPtr, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to create order",
			"details": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Order created successfully",
		"order":   order,
	})
}

// Get order by ID
func getOrderHandler(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	order, err := getOrderByID(orderID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Check if user has access to this order
	user := c.Locals("user")
	if user != nil {
		userClaims := user.(*Claims)
		// Allow access if it's the user's order or if user is admin
		if userClaims.Role != "admin" && (order.UserID == nil || *order.UserID != userClaims.UserID) {
			return c.Status(403).JSON(fiber.Map{
				"error": "Access denied",
			})
		}
	} else {
		// For guest users, check session ID
		sessionID := c.Get("X-Session-ID", "")
		if order.SessionID == nil || *order.SessionID != sessionID {
			return c.Status(403).JSON(fiber.Map{
				"error": "Access denied",
			})
		}
	}

	return c.JSON(order)
}

// Get user orders
func getUserOrdersHandler(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authentication required",
		})
	}

	userClaims := user.(*Claims)
	orders, err := getUserOrders(userClaims.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get orders",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"total":  len(orders),
	})
}

// Admin: Get all orders
func adminGetOrdersHandler(c *fiber.Ctx) error {
	orders, err := getAllOrders()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to get orders",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"total":  len(orders),
	})
}

// Admin: Update order status
func adminUpdateOrderStatusHandler(c *fiber.Ctx) error {
	orderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid order ID",
		})
	}

	var req UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err = updateOrderStatus(orderID, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to update order status",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Order status updated successfully",
	})
}

// =====================================================
// VALIDATION MIDDLEWARE
// =====================================================

// Validate email format
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// Validate input middleware
func validateRegistrationInput(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err == nil {
		// Store parsed request for the handler
		c.Locals("parsedRequest", req)

		// Validate email
		if !isValidEmail(req.Email) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid email format",
			})
		}

		// Validate password strength
		if len(req.Password) < 6 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Password must be at least 6 characters long",
			})
		}

		// Validate required fields
		if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Email) == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Username and email are required",
			})
		}
	}

	return c.Next()
}

// Validate product input
func validateProductInput(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err == nil {
		c.Locals("parsedRequest", req)

		if strings.TrimSpace(req.Name) == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Product name is required",
			})
		}

		if req.BasePrice <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Product price must be greater than 0",
			})
		}

		if req.CategoryID <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Valid category ID is required",
			})
		}
	}

	return c.Next()
}

// Validate cart input
func validateCartInput(c *fiber.Ctx) error {
	var req AddToCartRequest
	if err := c.BodyParser(&req); err == nil {
		c.Locals("parsedRequest", req)

		if req.ProductID <= 0 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Valid product ID is required",
			})
		}

		if req.Quantity <= 0 || req.Quantity > 100 {
			return c.Status(400).JSON(fiber.Map{
				"error": "Quantity must be between 1 and 100",
			})
		}
	}

	return c.Next()
}

// =====================================================
// LOGGING MIDDLEWARE
// =====================================================

// Request logging middleware
func loggingMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	// Log request details
	duration := time.Since(start)
	status := c.Response().StatusCode()

	log.Printf(
		"%s %s - %d - %v - %s",
		c.Method(),
		c.Path(),
		status,
		duration,
		c.IP(),
	)

	return err
}

// Error logging middleware
func errorLoggingMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		log.Printf("ERROR: %s %s - %v", c.Method(), c.Path(), err)
	}

	return err
}
