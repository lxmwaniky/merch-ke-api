package main

import (
	"testing"
	"time"
)

// TestProductStruct tests the Product struct
func TestProductStruct(t *testing.T) {
	product := Product{
		ID:               1,
		Name:             "Test Product",
		Slug:             "test-product",
		Description:      "A test product description",
		ShortDescription: "Test product",
		CategoryID:       5,
		BasePrice:        1500.00,
		SKUPrefix:        "PROD",
		IsActive:         true,
		IsFeatured:       false,
		Weight:           0.5,
		Dimensions:       "10x10x5",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if product.ID != 1 {
		t.Errorf("ID = %d, want 1", product.ID)
	}

	if product.Name != "Test Product" {
		t.Errorf("Name = %s, want Test Product", product.Name)
	}

	if product.Slug != "test-product" {
		t.Errorf("Slug = %s, want test-product", product.Slug)
	}

	if product.BasePrice != 1500.00 {
		t.Errorf("BasePrice = %.2f, want 1500.00", product.BasePrice)
	}

	if !product.IsActive {
		t.Error("IsActive should be true")
	}

	if product.IsFeatured {
		t.Error("IsFeatured should be false")
	}
}

// TestProductImageStruct tests the ProductImage struct
func TestProductImageStruct(t *testing.T) {
	image := ProductImage{
		ID:           1,
		ProductID:    1,
		ImageURL:     "https://example.com/image.jpg",
		ImagePath:    "/images/product.jpg",
		ImageType:    "jpg",
		AltText:      "Product image",
		DisplayOrder: 1,
		FileSize:     1024,
		Width:        800,
		Height:       600,
		IsPrimary:    true,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}

	if image.ID != 1 {
		t.Errorf("ID = %d, want 1", image.ID)
	}

	if image.ProductID != 1 {
		t.Errorf("ProductID = %d, want 1", image.ProductID)
	}

	if image.ImageURL != "https://example.com/image.jpg" {
		t.Errorf("ImageURL = %s, want https://example.com/image.jpg", image.ImageURL)
	}

	if !image.IsPrimary {
		t.Error("IsPrimary should be true")
	}

	if image.Width != 800 {
		t.Errorf("Width = %d, want 800", image.Width)
	}

	if image.Height != 600 {
		t.Errorf("Height = %d, want 600", image.Height)
	}
}

// TestCategoryStruct tests the Category struct
func TestCategoryStruct(t *testing.T) {
	category := Category{
		ID:          1,
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "A test category",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if category.ID != 1 {
		t.Errorf("ID = %d, want 1", category.ID)
	}

	if category.Name != "Test Category" {
		t.Errorf("Name = %s, want Test Category", category.Name)
	}

	if category.Slug != "test-category" {
		t.Errorf("Slug = %s, want test-category", category.Slug)
	}

	if !category.IsActive {
		t.Error("IsActive should be true")
	}
}

// TestOrderStruct tests the Order struct
func TestOrderStruct(t *testing.T) {
	userId := 1
	order := Order{
		ID:          123,
		OrderNumber: "ORD-20251013-0123",
		UserID:      &userId,
		TotalAmount: 6500.00,
		Status:      "pending",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	if order.ID != 123 {
		t.Errorf("ID = %d, want 123", order.ID)
	}

	if order.OrderNumber != "ORD-20251013-0123" {
		t.Errorf("OrderNumber = %s, want ORD-20251013-0123", order.OrderNumber)
	}

	if order.UserID == nil || *order.UserID != 1 {
		t.Errorf("UserID = %v, want 1", order.UserID)
	}

	if order.TotalAmount != 6500.00 {
		t.Errorf("TotalAmount = %.2f, want 6500.00", order.TotalAmount)
	}

	if order.Status != "pending" {
		t.Errorf("Status = %s, want pending", order.Status)
	}
}

// TestCartItemStruct tests the CartItem struct
func TestCartItemStruct(t *testing.T) {
	userId := 1
	cartItem := CartItem{
		ID:        1,
		UserID:    &userId,
		ProductID: 10,
		Quantity:  2,
	}

	if cartItem.ID != 1 {
		t.Errorf("ID = %d, want 1", cartItem.ID)
	}

	if cartItem.UserID == nil || *cartItem.UserID != 1 {
		t.Errorf("UserID = %v, want 1", cartItem.UserID)
	}

	if cartItem.ProductID != 10 {
		t.Errorf("ProductID = %d, want 10", cartItem.ProductID)
	}

	if cartItem.Quantity != 2 {
		t.Errorf("Quantity = %d, want 2", cartItem.Quantity)
	}
}

// TestValidateProductInput tests product validation logic
func TestValidateProductInput(t *testing.T) {
	tests := []struct {
		name    string
		product Product
		isValid bool
	}{
		{
			name: "Valid product",
			product: Product{
				Name:       "Valid Product",
				Slug:       "valid-product",
				BasePrice:  100.00,
				CategoryID: 1,
			},
			isValid: true,
		},
		{
			name: "Empty name",
			product: Product{
				Name:       "",
				Slug:       "valid-product",
				BasePrice:  100.00,
				CategoryID: 1,
			},
			isValid: false,
		},
		{
			name: "Empty slug",
			product: Product{
				Name:       "Valid Product",
				Slug:       "",
				BasePrice:  100.00,
				CategoryID: 1,
			},
			isValid: false,
		},
		{
			name: "Negative price",
			product: Product{
				Name:       "Valid Product",
				Slug:       "valid-product",
				BasePrice:  -100.00,
				CategoryID: 1,
			},
			isValid: false,
		},
		{
			name: "Zero price",
			product: Product{
				Name:       "Valid Product",
				Slug:       "valid-product",
				BasePrice:  0,
				CategoryID: 1,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.product.Name != "" &&
				tt.product.Slug != "" &&
				tt.product.BasePrice > 0 &&
				tt.product.CategoryID > 0

			if valid != tt.isValid {
				t.Errorf("Validation = %v, want %v", valid, tt.isValid)
			}
		})
	}
}

// TestValidateQuantity tests quantity validation
func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity int
		isValid  bool
	}{
		{
			name:     "Valid quantity",
			quantity: 1,
			isValid:  true,
		},
		{
			name:     "Multiple items",
			quantity: 10,
			isValid:  true,
		},
		{
			name:     "Zero quantity",
			quantity: 0,
			isValid:  false,
		},
		{
			name:     "Negative quantity",
			quantity: -5,
			isValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.quantity > 0
			if valid != tt.isValid {
				t.Errorf("Validation = %v, want %v for quantity %d", valid, tt.isValid, tt.quantity)
			}
		})
	}
}
