package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// Product struct to match database
type Product struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description"`
	ShortDescription string    `json:"short_description"`
	CategoryID       int       `json:"category_id"`
	BasePrice        float64   `json:"base_price"`
	SKUPrefix        string    `json:"sku_prefix"`
	ImageURL         string    `json:"image_url,omitempty"`
	IsActive         bool      `json:"is_active"`
	IsFeatured       bool      `json:"is_featured"`
	Weight           float64   `json:"weight"`
	Dimensions       string    `json:"dimensions"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// ProductImage struct for product images
type ProductImage struct {
	ID           int    `json:"id"`
	ProductID    int    `json:"product_id"`
	VariantID    *int   `json:"variant_id"`
	ImageURL     string `json:"image_url"`
	ImagePath    string `json:"image_path"`
	ImageType    string `json:"image_type"`
	AltText      string `json:"alt_text"`
	DisplayOrder int    `json:"display_order"`
	FileSize     int    `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsPrimary    bool   `json:"is_primary"`
	CreatedAt    string `json:"created_at"`
}

// ProductImageRequest struct for creating images
type ProductImageRequest struct {
	ImageURL     string `json:"image_url"`
	ImagePath    string `json:"image_path,omitempty"`
	ImageType    string `json:"image_type,omitempty"`
	AltText      string `json:"alt_text,omitempty"`
	DisplayOrder int    `json:"display_order,omitempty"`
	IsPrimary    bool   `json:"is_primary,omitempty"`
}

// CreateProductRequest struct for admin product creation
type CreateProductRequest struct {
	Name             string                `json:"name"`
	Slug             string                `json:"slug"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"short_description"`
	CategoryID       int                   `json:"category_id"`
	BasePrice        float64               `json:"base_price"`
	SKUPrefix        string                `json:"sku_prefix"`
	ImageURL         string                `json:"image_url,omitempty"`
	IsFeatured       bool                  `json:"is_featured"`
	Weight           float64               `json:"weight"`
	Dimensions       string                `json:"dimensions"`
	Images           []ProductImageRequest `json:"images,omitempty"`
}

// Category struct to match database
type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ParentID    *int      `json:"parent_id"`
	ImageURL    string    `json:"image_url"`
	IsActive    bool      `json:"is_active"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Get all products from database
func getProductsFromDB() ([]Product, error) {
	query := `
		SELECT p.id, p.name, p.slug, p.description, p.category_id, p.base_price, p.is_active, p.is_featured,
		       COALESCE((SELECT image_url FROM catalog.product_images WHERE product_id = p.id ORDER BY is_primary DESC, display_order LIMIT 1), '') as image_url
		FROM catalog.products p
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
		err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID, &p.BasePrice, &p.IsActive, &p.IsFeatured, &p.ImageURL)
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
		SELECT id, name, slug, description, parent_id, image_url, is_active, sort_order, created_at, updated_at
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
		err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.ParentID, &c.ImageURL, &c.IsActive, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
} // CreateCategoryRequest struct for admin category creation
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	ParentID    *int   `json:"parent_id,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateCategoryRequest struct for category updates
type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Slug        *string `json:"slug,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *int    `json:"parent_id,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	SortOrder   *int    `json:"sort_order,omitempty"`
}

// Create new category (admin only)
func createCategory(req *CreateCategoryRequest) (*Category, error) {
	query := `
		INSERT INTO catalog.categories (name, slug, description, parent_id, image_url, sort_order)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, slug, description, parent_id, image_url, is_active, sort_order, created_at, updated_at
	`

	var category Category
	err := db.QueryRow(query,
		req.Name, req.Slug, req.Description, req.ParentID, req.ImageURL, req.SortOrder,
	).Scan(&category.ID, &category.Name, &category.Slug, &category.Description,
		&category.ParentID, &category.ImageURL, &category.IsActive, &category.SortOrder,
		&category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Update existing category (admin only)
func updateCategory(id int, req *UpdateCategoryRequest) (*Category, error) {
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
	if req.ParentID != nil {
		setParts = append(setParts, fmt.Sprintf("parent_id = $%d", argIndex))
		args = append(args, *req.ParentID)
		argIndex++
	}
	if req.ImageURL != nil {
		setParts = append(setParts, fmt.Sprintf("image_url = $%d", argIndex))
		args = append(args, *req.ImageURL)
		argIndex++
	}
	if req.IsActive != nil {
		setParts = append(setParts, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}
	if req.SortOrder != nil {
		setParts = append(setParts, fmt.Sprintf("sort_order = $%d", argIndex))
		args = append(args, *req.SortOrder)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Add updated_at
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add ID for WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE catalog.categories 
		SET %s 
		WHERE id = $%d 
		RETURNING id, name, slug, description, parent_id, image_url, is_active, sort_order, created_at, updated_at
	`, strings.Join(setParts, ", "), argIndex)

	var category Category
	err := db.QueryRow(query, args...).Scan(
		&category.ID, &category.Name, &category.Slug, &category.Description,
		&category.ParentID, &category.ImageURL, &category.IsActive, &category.SortOrder,
		&category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

// Delete category (admin only) - soft delete
func deleteCategory(id int) error {
	query := `DELETE FROM catalog.categories WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// =====================================================
// PRODUCT IMAGE FUNCTIONS
// =====================================================

// Create product image
func createProductImage(productID int, variantID *int, req *ProductImageRequest) (*ProductImage, error) {
	query := `
		INSERT INTO catalog.product_images (product_id, variant_id, image_url, image_path, image_type, alt_text, display_order, is_primary)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, product_id, variant_id, image_url, image_path, image_type, alt_text, display_order, file_size, width, height, is_primary, created_at
	`

	var image ProductImage
	err := db.QueryRow(query, productID, variantID, req.ImageURL, req.ImagePath, req.ImageType, req.AltText, req.DisplayOrder, req.IsPrimary).Scan(
		&image.ID, &image.ProductID, &image.VariantID, &image.ImageURL, &image.ImagePath, &image.ImageType,
		&image.AltText, &image.DisplayOrder, &image.FileSize, &image.Width, &image.Height, &image.IsPrimary, &image.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

// Get product images by product ID
func getProductImages(productID int) ([]ProductImage, error) {
	query := `
		SELECT id, product_id, variant_id, image_url, image_path, image_type, alt_text, display_order, file_size, width, height, is_primary, created_at
		FROM catalog.product_images 
		WHERE product_id = $1 
		ORDER BY display_order ASC, is_primary DESC, created_at ASC
	`

	rows, err := db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []ProductImage
	for rows.Next() {
		var img ProductImage
		err := rows.Scan(&img.ID, &img.ProductID, &img.VariantID, &img.ImageURL, &img.ImagePath, &img.ImageType,
			&img.AltText, &img.DisplayOrder, &img.FileSize, &img.Width, &img.Height, &img.IsPrimary, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

// Update product image
func updateProductImage(id int, req *ProductImageRequest) (*ProductImage, error) {
	query := `
		UPDATE catalog.product_images 
		SET image_url = $2, image_path = $3, image_type = $4, alt_text = $5, display_order = $6, is_primary = $7
		WHERE id = $1
		RETURNING id, product_id, variant_id, image_url, image_path, image_type, alt_text, display_order, file_size, width, height, is_primary, created_at
	`

	var image ProductImage
	err := db.QueryRow(query, id, req.ImageURL, req.ImagePath, req.ImageType, req.AltText, req.DisplayOrder, req.IsPrimary).Scan(
		&image.ID, &image.ProductID, &image.VariantID, &image.ImageURL, &image.ImagePath, &image.ImageType,
		&image.AltText, &image.DisplayOrder, &image.FileSize, &image.Width, &image.Height, &image.IsPrimary, &image.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

// Delete product image
func deleteProductImage(id int) error {
	query := `DELETE FROM catalog.product_images WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
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
	ImageURL         *string  `json:"image_url,omitempty"`
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

	// Handle image_url update if provided
	if req.ImageURL != nil && *req.ImageURL != "" {
		fmt.Printf("🖼️  Processing image_url update for product %d: %s\n", id, *req.ImageURL)
		
		// Check if product already has an image
		var existingImageID int
		checkQuery := `SELECT id FROM catalog.product_images WHERE product_id = $1 AND is_primary = true LIMIT 1`
		err := db.QueryRow(checkQuery, id).Scan(&existingImageID)
		
		if err == sql.ErrNoRows {
			// No existing image, create new one
			fmt.Printf("📸 No existing image found, creating new one for product %d\n", id)
			insertQuery := `
				INSERT INTO catalog.product_images (product_id, image_url, image_path, is_primary, display_order)
				VALUES ($1, $2, $2, true, 0)
			`
			_, err = db.Exec(insertQuery, id, *req.ImageURL)
			if err != nil {
				fmt.Printf("❌ Failed to insert product image: %v\n", err)
			} else {
				fmt.Printf("✅ Successfully created new product image for product %d\n", id)
			}
		} else if err == nil {
			// Update existing image
			fmt.Printf("📸 Found existing image ID %d, updating...\n", existingImageID)
			updateQuery := `UPDATE catalog.product_images SET image_url = $1, image_path = $1 WHERE id = $2`
			_, err = db.Exec(updateQuery, *req.ImageURL, existingImageID)
			if err != nil {
				fmt.Printf("❌ Failed to update product image: %v\n", err)
			} else {
				fmt.Printf("✅ Successfully updated product image ID %d with URL: %s\n", existingImageID, *req.ImageURL)
			}
		} else {
			// Some other error occurred
			fmt.Printf("❌ Error checking for existing image: %v\n", err)
		}
	} else {
		if req.ImageURL == nil {
			fmt.Printf("ℹ️  No image_url provided in update request for product %d\n", id)
		} else {
			fmt.Printf("ℹ️  Empty image_url provided in update request for product %d\n", id)
		}
	}

	return &product, nil
}

// Hard delete product (admin only) - with validation
func deleteProduct(id int) error {
	// First, check if product has any variants that were used in orders
	var orderCount int
	checkQuery := `
		SELECT COUNT(DISTINCT oi.order_id)
		FROM orders.order_items oi
		JOIN catalog.product_variants pv ON oi.variant_id = pv.id
		WHERE pv.product_id = $1
	`
	err := db.QueryRow(checkQuery, id).Scan(&orderCount)
	if err != nil {
		return fmt.Errorf("failed to check product orders: %v", err)
	}

	if orderCount > 0 {
		return fmt.Errorf("cannot delete product: it has been used in %d order(s). Consider marking it as inactive instead", orderCount)
	}

	// If no orders, proceed with hard delete (CASCADE will handle variants and images)
	query := `DELETE FROM catalog.products WHERE id = $1`

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
	VariantID   *int    `json:"variant_id,omitempty"`
	ProductName string  `json:"product_name"`
	VariantSKU  string  `json:"variant_sku"`
	Size        *string `json:"size,omitempty"`
	Color       *string `json:"color,omitempty"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

// CreateOrderRequest represents order creation request
type CreateOrderRequest struct {
	ShippingAddress *string `json:"shipping_address,omitempty"`
	BillingAddress  *string `json:"billing_address,omitempty"`
	PaymentMethod   *string `json:"payment_method,omitempty"`
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
	log.Printf("🔵 createOrderFromCart STARTED - userID=%v, sessionID=%v", userID, sessionID)

	// Generate unique order number
	orderNumber := generateOrderNumber()
	log.Printf("🔵 Generated order number: %s", orderNumber)

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("❌ Failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			log.Printf("🔴 Rolling back transaction")
			tx.Rollback()
		} else {
			log.Printf("🟢 Committing transaction")
			tx.Commit()
		}
	}()

	// Get cart items
	log.Printf("🔵 Fetching cart items...")
	var cartItems []CartItem
	if userID != nil {
		cartItems, err = getUserCartItems(*userID)
	} else if sessionID != nil {
		cartItems, err = getGuestCartItems(*sessionID)
	}

	if err != nil {
		log.Printf("❌ Failed to get cart items: %v", err)
		return nil, err
	}

	if len(cartItems) == 0 {
		log.Printf("❌ Cart is empty")
		return nil, fmt.Errorf("cart is empty")
	}

	log.Printf("🟢 Found %d cart items", len(cartItems))

	// Calculate total amount
	var totalAmount float64
	for _, item := range cartItems {
		totalAmount += float64(item.Quantity) * item.Price
	}

	// Create order
	var orderID int
	orderQuery := `
		INSERT INTO orders.orders (
			user_id, session_id, order_number, status, 
			subtotal, total_amount, payment_status, payment_method,
			shipping_address, billing_address, notes
		)
		VALUES ($1, $2, $3, 'pending', $4, $5, 'pending', $6, $7, $8, $9)
		RETURNING id
	`

	err = tx.QueryRow(
		orderQuery,
		userID, sessionID, orderNumber,
		totalAmount, totalAmount, req.PaymentMethod, // subtotal = total for now
		req.ShippingAddress, req.BillingAddress, req.Notes,
	).Scan(&orderID)
	if err != nil {
		log.Printf("❌ ORDER INSERT FAILED - SQL Error: %v", err)
		log.Printf("   Values: userID=%v, sessionID=%v, orderNumber=%s", userID, sessionID, orderNumber)
		log.Printf("   Amounts: subtotal=%f, total=%f", totalAmount, totalAmount)
		log.Printf("   Payment: method=%v, status=pending", req.PaymentMethod)
		log.Printf("   Addresses: shipping=%v, billing=%v", req.ShippingAddress, req.BillingAddress)
		log.Printf("   Notes: %v", req.Notes)
		return nil, err
	}

	// Create order items
	for _, item := range cartItems {
		totalPrice := float64(item.Quantity) * item.Price

		// For now, use product name and generate a simple SKU
		// TODO: Implement proper variant support
		productName := item.ProductName
		if productName == "" {
			productName = "Product"
		}
		variantSKU := fmt.Sprintf("PROD-%d", item.ProductID)

		orderItemQuery := `
			INSERT INTO orders.order_items (
				order_id, product_name, variant_sku, 
				unit_price, quantity, total_price
			)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		log.Printf("🔵 Inserting order item: %s (SKU: %s, Price: %.2f, Qty: %d)", productName, variantSKU, item.Price, item.Quantity)
		_, err = tx.Exec(orderItemQuery, orderID, productName, variantSKU, item.Price, item.Quantity, totalPrice)
		if err != nil {
			log.Printf("❌ Failed to insert order item: %v", err)
			return nil, err
		}
	}

	// Clear cart after order creation
	log.Printf("🔵 Clearing cart...")
	if userID != nil {
		_, err = tx.Exec("DELETE FROM orders.cart_items WHERE user_id = $1", *userID)
	} else if sessionID != nil {
		_, err = tx.Exec("DELETE FROM orders.guest_cart_items WHERE session_id = $1", *sessionID)
	}

	if err != nil {
		log.Printf("❌ Failed to clear cart: %v", err)
		return nil, err
	}
	log.Printf("🟢 Cart cleared")

	// Get the created order using the transaction
	log.Printf("🔵 Fetching created order (ID: %d) using transaction...", orderID)
	query := `
		SELECT id, user_id, session_id, order_number, status, total_amount, payment_status, 
		       payment_method, shipping_address, billing_address, notes, created_at, updated_at
		FROM orders.orders 
		WHERE id = $1
	`

	var order Order
	err = tx.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.SessionID, &order.OrderNumber,
		&order.Status, &order.TotalAmount, &order.PaymentStatus,
		&order.PaymentMethod, &order.ShippingAddress, &order.BillingAddress,
		&order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		log.Printf("❌ Failed to fetch order from transaction: %v", err)
		return nil, err
	}

	log.Printf("🟢 Order fetched successfully")
	return &order, nil
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
		SELECT id, order_id, product_name, variant_sku, size, color,
		       unit_price, quantity, total_price
		FROM orders.order_items
		WHERE order_id = $1
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
			&item.ID, &item.OrderID, &item.ProductName, &item.VariantSKU,
			&item.Size, &item.Color, &item.UnitPrice, &item.Quantity, &item.TotalPrice,
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
