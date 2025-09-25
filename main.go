package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Initialize database connection
	initDatabase()
	defer closeDatabase()

	app := fiber.New(fiber.Config{
		AppName: "Merch Ke API",
	})

	app.Use(cors.New())
	app.Use(loggingMiddleware)
	app.Use(errorLoggingMiddleware)

	// Public routes
	app.Get("/health", healthHandler)
	app.Get("/api/products", productsHandler)
	app.Get("/api/products/:id", singleProductHandler)
	app.Get("/api/categories", categoriesHandler)

	// Authentication routes
	app.Post("/api/auth/register", registerHandler)
	app.Post("/api/auth/login", loginHandler)

	// Protected routes (require authentication)
	app.Get("/api/auth/profile", authMiddleware, profileHandler)

	// Cart routes (work for both authenticated and guest users)
	app.Post("/api/cart", optionalAuthMiddleware, addToCartHandler)
	app.Get("/api/cart", optionalAuthMiddleware, getCartHandler)
	app.Put("/api/cart/:productId", optionalAuthMiddleware, updateCartHandler)
	app.Delete("/api/cart/:productId", optionalAuthMiddleware, removeFromCartHandler)

	// Cart migration route (for when guest users register/login)
	app.Post("/api/cart/migrate", authMiddleware, migrateCartHandler)

	// Points routes (authenticated users only)
	app.Get("/api/points", authMiddleware, getUserPointsHandler)

	// Order routes
	app.Post("/api/orders", optionalAuthMiddleware, createOrderHandler)           // Create order from cart
	app.Get("/api/orders/:id", optionalAuthMiddleware, getOrderHandler)          // Get specific order
	app.Get("/api/orders", authMiddleware, getUserOrdersHandler)                 // Get user's orders

	// Admin routes (require admin privileges)
	admin := app.Group("/api/admin", authMiddleware, adminMiddleware)
	admin.Post("/products", adminCreateProductHandler)
	admin.Put("/products/:id", adminUpdateProductHandler)
	admin.Delete("/products/:id", adminDeleteProductHandler)
	admin.Get("/products", adminGetProductsHandler)
	admin.Get("/orders", adminGetOrdersHandler)                    // Get all orders
	admin.Put("/orders/:id/status", adminUpdateOrderStatusHandler) // Update order status
	admin.Post("/init-cart-tables", initCartTablesHandler)         // Temporary endpoint
	admin.Post("/init-order-tables", initOrderTablesHandler)       // Initialize order tables
	admin.Post("/fix-cart-tables", fixCartTablesHandler)           // Fix cart tables structure
	admin.Post("/init-order-tables", initOrderTablesHandler)       // Create order tables

	// Temporary public endpoint to promote user (REMOVE IN PRODUCTION)
	app.Post("/api/temp/promote-admin", promoteToAdminHandler)

	log.Println("🚀 Merch Ke API starting on http://161.35.104.94:8080")
	log.Fatal(app.Listen("0.0.0.0:8080"))
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "Merch Ke API",
	})
}

func productsHandler(c *fiber.Ctx) error {
	products, err := getProductsFromDB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    len(products),
	})
}

func singleProductHandler(c *fiber.Ctx) error {
	idStr := c.Params("id")

	// Convert string ID to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	// Get product from database
	product, err := getProductByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return c.Status(404).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch product",
			"details": err.Error(),
		})
	}

	return c.JSON(product)
}

func categoriesHandler(c *fiber.Ctx) error {
	categories, err := getCategoriesFromDB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch categories",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"categories": categories,
		"total":      len(categories),
	})
}

func parseID(id string) int {
	if id == "1" {
		return 1
	}
	if id == "2" {
		return 2
	}
	return 0
}
