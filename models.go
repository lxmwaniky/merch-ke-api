package main

import (
	"fmt"
	"strings"
	"time"
)

// Product struct to match database
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	CategoryID  int     `json:"category_id"`
	BasePrice   float64 `json:"base_price"`
	IsActive    bool    `json:"is_active"`
	IsFeatured  bool    `json:"is_featured"`
}

// Category struct to match database
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id"`
	IsActive    bool   `json:"is_active"`
}

// Get all products from database
func getProductsFromDB() ([]Product, error) {
	query := `
		SELECT id, name, slug, description, category_id, base_price, is_active, is_featured 
		FROM catalog.products 
		WHERE is_active = true 
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID, &p.BasePrice, &p.IsActive, &p.IsFeatured)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// Get single product by ID
func getProductByID(id int) (*Product, error) {
	query := `
		SELECT id, name, slug, description, category_id, base_price, is_active, is_featured 
		FROM catalog.products 
		WHERE id = $1 AND is_active = true
	`

	var p Product
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID, &p.BasePrice, &p.IsActive, &p.IsFeatured)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Get all categories from database
func getCategoriesFromDB() ([]Category, error) {
	query := `
		SELECT id, name, slug, description, parent_id, is_active 
		FROM catalog.categories 
		WHERE is_active = true 
		ORDER BY sort_order, name
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.IsActive)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

// CreateProductRequest struct for admin product creation
type CreateProductRequest struct {
	Name             string  `json:"name"`
	Slug             string  `json:"slug"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	CategoryID       int     `json:"category_id"`
	BasePrice        float64 `json:"base_price"`
	SKUPrefix        string  `json:"sku_prefix"`
	IsFeatured       bool    `json:"is_featured"`
	Weight           float64 `json:"weight"`
	Dimensions       string  `json:"dimensions"`
}

// UpdateProductRequest struct for product updates
type UpdateProductRequest struct {
	Name             *string  `json:"name,omitempty"`
	Slug             *string  `json:"slug,omitempty"`
	Description      *string  `json:"description,omitempty"`
	ShortDescription *string  `json:"short_description,omitempty"`
	CategoryID       *int     `json:"category_id,omitempty"`
	BasePrice        *float64 `json:"base_price,omitempty"`
	SKUPrefix        *string  `json:"sku_prefix,omitempty"`
	IsFeatured       *bool    `json:"is_featured,omitempty"`
	IsActive         *bool    `json:"is_active,omitempty"`
	Weight           *float64 `json:"weight,omitempty"`
	Dimensions       *string  `json:"dimensions,omitempty"`
}

// Create new product (admin only)
func createProduct(req *CreateProductRequest) (*Product, error) {
	query := `
		INSERT INTO catalog.products (name, slug, description, short_description, category_id, base_price, sku_prefix, is_featured, weight, dimensions)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, slug, description, category_id, base_price, is_active, is_featured
	`

	var product Product
	err := db.QueryRow(query,
		req.Name, req.Slug, req.Description, req.ShortDescription,
		req.CategoryID, req.BasePrice, req.SKUPrefix, req.IsFeatured,
		req.Weight, req.Dimensions,
	).Scan(&product.ID, &product.Name, &product.Slug, &product.Description,
		&product.CategoryID, &product.BasePrice, &product.IsActive, &product.IsFeatured)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Update existing product (admin only)
func updateProduct(id int, req *UpdateProductRequest) (*Product, error) {
	// Build dynamic query based on provided fields
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Slug != nil {
		setParts = append(setParts, fmt.Sprintf("slug = $%d", argIndex))
		args = append(args, *req.Slug)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.ShortDescription != nil {
		setParts = append(setParts, fmt.Sprintf("short_description = $%d", argIndex))
		args = append(args, *req.ShortDescription)
		argIndex++
	}
	if req.CategoryID != nil {
		setParts = append(setParts, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.BasePrice != nil {
		setParts = append(setParts, fmt.Sprintf("base_price = $%d", argIndex))
		args = append(args, *req.BasePrice)
		argIndex++
	}
	if req.IsFeatured != nil {
		setParts = append(setParts, fmt.Sprintf("is_featured = $%d", argIndex))
		args = append(args, *req.IsFeatured)
		argIndex++
	}
	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add updated_at timestamp
	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))

	// Add product ID for WHERE clause
	args = append(args, id)
	whereClause := fmt.Sprintf("WHERE id = $%d", argIndex)

	query := fmt.Sprintf(`
		UPDATE catalog.products 
		SET %s 
		%s
		RETURNING id, name, slug, description, category_id, base_price, is_active, is_featured
	`, strings.Join(setParts, ", "), whereClause)

	var product Product
	err := db.QueryRow(query, args...).Scan(
		&product.ID, &product.Name, &product.Slug, &product.Description,
		&product.CategoryID, &product.BasePrice, &product.IsActive, &product.IsFeatured,
	)

	if err != nil {
		return nil, err
	}

	return &product, nil
}

// Soft delete product (admin only)
func deleteProduct(id int) error {
	query := `UPDATE catalog.products SET is_active = false, updated_at = NOW() WHERE id = $1`

	result, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

// =====================================================
// CART MODELS AND FUNCTIONS
// =====================================================

// CartItem represents a cart item
type CartItem struct {
	ID        int     `json:"id"`
	UserID    *int    `json:"user_id,omitempty"`
	SessionID *string `json:"session_id,omitempty"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	// Joined fields from product
	ProductName string  `json:"product_name"`
	ProductSlug string  `json:"product_slug"`
	Price       float64 `json:"price"`
	ImageURL    *string `json:"image_url,omitempty"`
}

// CartSummary represents cart totals
type CartSummary struct {
	Items      []CartItem `json:"items"`
	TotalItems int        `json:"total_items"`
	Subtotal   float64    `json:"subtotal"`
}

// AddToCartRequest represents add to cart request
type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// UserPoints represents user points balance
type UserPoints struct {
	UserID        int `json:"user_id"`
	PointsBalance int `json:"points_balance"`
	TotalEarned   int `json:"total_earned"`
	TotalSpent    int `json:"total_spent"`
}

// Order represents an order
type Order struct {
	ID              int         `json:"id"`
	UserID          *int        `json:"user_id,omitempty"` // nil for guest orders
	SessionID       *string     `json:"session_id,omitempty"`
	OrderNumber     string      `json:"order_number"`
	Status          string      `json:"status"` // pending, confirmed, processing, shipped, delivered, cancelled
	TotalAmount     float64     `json:"total_amount"`
	PaymentStatus   string      `json:"payment_status"` // pending, paid, failed, refunded
	PaymentMethod   *string     `json:"payment_method,omitempty"`
	ShippingAddress *string     `json:"shipping_address,omitempty"`
	BillingAddress  *string     `json:"billing_address,omitempty"`
	Notes           *string     `json:"notes,omitempty"`
	CreatedAt       string      `json:"created_at"`
	UpdatedAt       string      `json:"updated_at"`
	Items           []OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          int     `json:"id"`
	OrderID     int     `json:"order_id"`
	ProductID   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	ProductName string  `json:"product_name"`
	ProductSlug string  `json:"product_slug"`
}

// CreateOrderRequest represents order creation request
type CreateOrderRequest struct {
	ShippingAddress *string `json:"shipping_address,omitempty"`
	BillingAddress  *string `json:"billing_address,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

// UpdateOrderStatusRequest represents order status update
type UpdateOrderStatusRequest struct {
	Status        string  `json:"status"`
	PaymentStatus *string `json:"payment_status,omitempty"`
	PaymentMethod *string `json:"payment_method,omitempty"`
}

// Add item to user cart (authenticated users)
func addToUserCart(userID, productID, quantity int) error {
	query := `
		INSERT INTO orders.cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET 
			quantity = orders.cart_items.quantity + $3,
			updated_at = NOW()
	`

	_, err := db.Exec(query, userID, productID, quantity)
	return err
}

// Add item to guest cart (session-based)
func addToGuestCart(sessionID string, productID, quantity int) error {
	query := `
		INSERT INTO orders.guest_cart_items (session_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (session_id, product_id)
		DO UPDATE SET 
			quantity = orders.guest_cart_items.quantity + $3,
			updated_at = NOW()
	`

	_, err := db.Exec(query, sessionID, productID, quantity)
	return err
}

// Get user cart items
func getUserCartItems(userID int) ([]CartItem, error) {
	query := `
		SELECT 
			ci.id, ci.user_id, ci.product_id, ci.quantity,
			p.name as product_name, p.slug as product_slug,
			p.base_price as price
		FROM orders.cart_items ci
		JOIN catalog.products p ON ci.product_id = p.id
		WHERE ci.user_id = $1 AND p.is_active = true
		ORDER BY ci.created_at DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.Quantity,
			&item.ProductName, &item.ProductSlug, &item.Price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Get guest cart items
func getGuestCartItems(sessionID string) ([]CartItem, error) {
	query := `
		SELECT 
			gci.id, gci.product_id, gci.quantity,
			p.name as product_name, p.slug as product_slug,
			p.base_price as price
		FROM orders.guest_cart_items gci
		JOIN catalog.products p ON gci.product_id = p.id
		WHERE gci.session_id = $1 AND p.is_active = true
		ORDER BY gci.created_at DESC
	`

	rows, err := db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		sessionIDStr := sessionID
		item.SessionID = &sessionIDStr

		err := rows.Scan(
			&item.ID, &item.ProductID, &item.Quantity,
			&item.ProductName, &item.ProductSlug, &item.Price,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Update cart item quantity
func updateCartItemQuantity(userID *int, sessionID *string, productID, quantity int) error {
	if userID != nil {
		// User cart
		if quantity <= 0 {
			query := `DELETE FROM orders.cart_items WHERE user_id = $1 AND product_id = $2`
			_, err := db.Exec(query, *userID, productID)
			return err
		} else {
			query := `UPDATE orders.cart_items SET quantity = $3, updated_at = NOW() WHERE user_id = $1 AND product_id = $2`
			_, err := db.Exec(query, *userID, productID, quantity)
			return err
		}
	} else if sessionID != nil {
		// Guest cart
		if quantity <= 0 {
			query := `DELETE FROM orders.guest_cart_items WHERE session_id = $1 AND product_id = $2`
			_, err := db.Exec(query, *sessionID, productID)
			return err
		} else {
			query := `UPDATE orders.guest_cart_items SET quantity = $3, updated_at = NOW() WHERE session_id = $1 AND product_id = $2`
			_, err := db.Exec(query, *sessionID, productID, quantity)
			return err
		}
	}

	return fmt.Errorf("either userID or sessionID must be provided")
}

// Remove item from cart
func removeFromCart(userID *int, sessionID *string, productID int) error {
	if userID != nil {
		query := `DELETE FROM orders.cart_items WHERE user_id = $1 AND product_id = $2`
		_, err := db.Exec(query, *userID, productID)
		return err
	} else if sessionID != nil {
		query := `DELETE FROM orders.guest_cart_items WHERE session_id = $1 AND product_id = $2`
		_, err := db.Exec(query, *sessionID, productID)
		return err
	}

	return fmt.Errorf("either userID or sessionID must be provided")
}

// Migrate guest cart to user cart when user registers/logs in
func migrateGuestCartToUser(sessionID string, userID int) error {
	// First, get guest cart items
	guestItems, err := getGuestCartItems(sessionID)
	if err != nil {
		return err
	}

	// Add each item to user cart
	for _, item := range guestItems {
		err := addToUserCart(userID, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Clear guest cart
	query := `DELETE FROM orders.guest_cart_items WHERE session_id = $1`
	_, err = db.Exec(query, sessionID)
	return err
}

// Get cart summary with totals
func getCartSummary(userID *int, sessionID *string) (*CartSummary, error) {
	var items []CartItem
	var err error

	if userID != nil {
		items, err = getUserCartItems(*userID)
	} else if sessionID != nil {
		items, err = getGuestCartItems(*sessionID)
	} else {
		return nil, fmt.Errorf("either userID or sessionID must be provided")
	}

	if err != nil {
		return nil, err
	}

	// Calculate totals
	totalItems := 0
	subtotal := 0.0

	for _, item := range items {
		totalItems += item.Quantity
		subtotal += float64(item.Quantity) * item.Price
	}

	return &CartSummary{
		Items:      items,
		TotalItems: totalItems,
		Subtotal:   subtotal,
	}, nil
}

// Initialize user points when user registers
func initializeUserPoints(userID int) error {
	query := `
		INSERT INTO auth.user_points (user_id, points_balance, total_earned, total_spent)
		VALUES ($1, 0, 0, 0)
		ON CONFLICT (user_id) DO NOTHING
	`
	_, err := db.Exec(query, userID)
	return err
}

// Add points to user account
func addPointsToUser(userID int, points int, description string, orderID *int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update points balance
	query1 := `
		UPDATE auth.user_points 
		SET points_balance = points_balance + $2,
		    total_earned = total_earned + $2,
		    updated_at = NOW()
		WHERE user_id = $1
	`
	_, err = tx.Exec(query1, userID, points)
	if err != nil {
		return err
	}

	// Record transaction
	query2 := `
		INSERT INTO auth.points_transactions (user_id, order_id, transaction_type, points, description)
		VALUES ($1, $2, 'earned', $3, $4)
	`
	_, err = tx.Exec(query2, userID, orderID, points, description)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Get user points balance
func getUserPoints(userID int) (*UserPoints, error) {
	query := `
		SELECT user_id, points_balance, total_earned, total_spent
		FROM auth.user_points
		WHERE user_id = $1
	`

	var points UserPoints
	err := db.QueryRow(query, userID).Scan(
		&points.UserID, &points.PointsBalance,
		&points.TotalEarned, &points.TotalSpent,
	)

	if err != nil {
		return nil, err
	}

	return &points, nil
}

// =====================================================
// ORDER MANAGEMENT FUNCTIONS
// =====================================================

// Create order from cart
func createOrderFromCart(userID *int, sessionID *string, req *CreateOrderRequest) (*Order, error) {
	// Generate unique order number
	orderNumber := generateOrderNumber()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Get cart items
	var cartItems []CartItem
	if userID != nil {
		cartItems, err = getUserCartItems(*userID)
	} else if sessionID != nil {
		cartItems, err = getGuestCartItems(*sessionID)
	}

	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Price
	}

	// Create order
	var orderID int
	orderQuery := `
		INSERT INTO orders.orders (user_id, session_id, order_number, status, total_amount, payment_status, shipping_address, billing_address, notes)
		VALUES ($1, $2, $3, 'pending', $4, 'pending', $5, $6, $7)
		RETURNING id
	`

	err = tx.QueryRow(orderQuery, userID, sessionID, orderNumber, totalAmount, req.ShippingAddress, req.BillingAddress, req.Notes).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	// Create order items
	for _, item := range cartItems {
		totalPrice := float64(item.Quantity) * item.Price
		orderItemQuery := `
			INSERT INTO orders.order_items (order_id, product_id, quantity, unit_price, total_price)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.Exec(orderItemQuery, orderID, item.ProductID, item.Quantity, item.Price, totalPrice)
		if err != nil {
			return nil, err
		}
	}

	// Clear cart after order creation
	if userID != nil {
		_, err = tx.Exec("DELETE FROM orders.cart_items WHERE user_id = $1", *userID)
	} else if sessionID != nil {
		_, err = tx.Exec("DELETE FROM orders.guest_cart_items WHERE session_id = $1", *sessionID)
	}

	if err != nil {
		return nil, err
	}

	// Get the created order
	order, err := getOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// Get order by ID
func getOrderByID(orderID int) (*Order, error) {
	query := `
		SELECT id, user_id, session_id, order_number, status, total_amount, payment_status, 
		       payment_method, shipping_address, billing_address, notes, created_at, updated_at
		FROM orders.orders 
		WHERE id = $1
	`

	var order Order
	err := db.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.SessionID, &order.OrderNumber,
		&order.Status, &order.TotalAmount, &order.PaymentStatus,
		&order.PaymentMethod, &order.ShippingAddress, &order.BillingAddress,
		&order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.total_price,
		       p.name as product_name, p.slug as product_slug
		FROM orders.order_items oi
		JOIN catalog.products p ON oi.product_id = p.id
		WHERE oi.order_id = $1
	`

	rows, err := db.Query(itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity,
			&item.UnitPrice, &item.TotalPrice, &item.ProductName, &item.ProductSlug,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	return &order, nil
}

// Get user orders
func getUserOrders(userID int) ([]Order, error) {
	query := `
		SELECT id, user_id, session_id, order_number, status, total_amount, payment_status,
		       payment_method, shipping_address, billing_address, notes, created_at, updated_at
		FROM orders.orders 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.SessionID, &order.OrderNumber,
			&order.Status, &order.TotalAmount, &order.PaymentStatus,
			&order.PaymentMethod, &order.ShippingAddress, &order.BillingAddress,
			&order.Notes, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Update order status (admin function)
func updateOrderStatus(orderID int, req *UpdateOrderStatusRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Status != "" {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, req.Status)
		argIndex++
	}

	if req.PaymentStatus != nil {
		setParts = append(setParts, fmt.Sprintf("payment_status = $%d", argIndex))
		args = append(args, *req.PaymentStatus)
		argIndex++
	}

	if req.PaymentMethod != nil {
		setParts = append(setParts, fmt.Sprintf("payment_method = $%d", argIndex))
		args = append(args, *req.PaymentMethod)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setParts = append(setParts, fmt.Sprintf("updated_at = NOW()"))

	query := fmt.Sprintf("UPDATE orders.orders SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIndex)
	args = append(args, orderID)

	_, err := db.Exec(query, args...)
	return err
}

// Get all orders (admin function)
func getAllOrders() ([]Order, error) {
	query := `
		SELECT id, user_id, session_id, order_number, status, total_amount, payment_status,
		       payment_method, shipping_address, billing_address, notes, created_at, updated_at
		FROM orders.orders 
		ORDER BY created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.SessionID, &order.OrderNumber,
			&order.Status, &order.TotalAmount, &order.PaymentStatus,
			&order.PaymentMethod, &order.ShippingAddress, &order.BillingAddress,
			&order.Notes, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// Generate unique order number
func generateOrderNumber() string {
	return fmt.Sprintf("ORD-%d", time.Now().Unix())
}
