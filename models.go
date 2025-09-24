package main

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
