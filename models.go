package main

import (
	"fmt"
	"strings"
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
		FROM products 
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
		FROM products 
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
		FROM categories 
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
		INSERT INTO products (name, slug, description, short_description, category_id, base_price, sku_prefix, is_featured, weight, dimensions)
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
		UPDATE products 
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
	query := `UPDATE products SET is_active = false, updated_at = NOW() WHERE id = $1`
	
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
