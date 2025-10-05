# 🚀 Merch Ke API - Postman Testing Guide

## 📋 Prerequisites

### 1. **Database Setup**
Before testing, ensure your PostgreSQL database is running and properly set up:

```bash
# Create database and user (if not done already)
psql -U postgres
CREATE DATABASE "merch-ke-db";
CREATE USER "merch-ke-admin" WITH PASSWORD 'merch-ke-password';
GRANT ALL PRIVILEGES ON DATABASE "merch-ke-db" TO "merch-ke-admin";
```

Run the schema file (creates multi-schema architecture):
```bash
psql -U merch-ke-admin -d merch-ke-db -f database/schema.sql
```

**📊 Database Schema Structure:**
The API now uses a multi-schema architecture for better organization:
- **`auth`** schema: User authentication, addresses, points
- **`catalog`** schema: Products, categories, variants, images  
- **`orders`** schema: Orders, cart items, order items

### 2. **Environment Variables**
Make sure your `.env` file is configured properly:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=merch-ke-admin
DB_PASSWORD=merch-ke-password
DB_NAME=merch-ke-db
DB_SSLMODE=disable
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

### 3. **Start the Server**
```bash
go run .
```

Server should start on: `http://localhost:8080`

---

## 🗂️ **Quick Schema Reference**

**Database Tables by Schema:**

| Schema | Table | Purpose |
|--------|--------|---------|
| `auth` | `auth.users` | User accounts & authentication |
| `auth` | `auth.user_addresses` | User shipping addresses |
| `auth` | `auth.user_points` | Loyalty points balance |
| `auth` | `auth.points_transactions` | Points transaction history |
| `catalog` | `catalog.products` | Product catalog |
| `catalog` | `catalog.categories` | Product categories |
| `catalog` | `catalog.product_variants` | Product sizes/colors/SKUs |
| `catalog` | `catalog.product_images` | Product image URLs |
| `orders` | `orders.orders` | Customer orders |
| `orders` | `orders.order_items` | Individual items in orders |
| `orders` | `orders.cart_items` | User shopping cart |
| `orders` | `orders.guest_cart_items` | Guest user cart |

**⚠️ Important:** Always use schema prefixes when working directly with the database!

---

## ✅ **Testing Progress Tracker**

**How to use:** Mark checkboxes as you complete each test case:
- `- [ ]` = Not tested yet
- `- [x]` = Test completed (click the checkbox in preview mode)

**Legend:**
- 🟢 = Success test (should return 200/201)  
- 🔴 = Error test (should return 4xx/5xx)

---

## 🧪 Testing Scenarios & Test Cases

### **Phase 1: Basic API Health Check**

#### 🟢 **Test 1: Health Check** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/health`
- **Expected Response:**
```json
{
  "status": "healthy",
  "service": "Merch Ke API"
}
```

---

### **Phase 2: Authentication Testing**

#### 🟢 **Test 2: User Registration** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/auth/register`
- **Headers:**
```
Content-Type: application/json
```
- **Body (JSON):**
```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123",
  "first_name": "Test",
  "last_name": "User",
  "phone": "+254700000000"
}
```
- **Expected Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "first_name": "Test",
    "last_name": "User",
    "phone": "+254700000000",
    "role": "customer",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-10-02T20:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 🟢 **Test 3: User Login** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/auth/login`
- **Headers:**
```
Content-Type: application/json
```
- **Body (JSON):**
```json
{
  "email": "test@example.com",
  "password": "password123"
}
```
- **Expected Response (200):**
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "first_name": "Test",
    "last_name": "User",
    "phone": "+254700000000",
    "role": "customer",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-10-02T20:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**💾 Save the JWT token for authenticated requests!**

#### 🟢 **Test 4: Get User Profile (Protected Route)** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/auth/profile`
- **Headers:**
```
Authorization: Bearer <your-jwt-token-here>
```
- **Expected Response (200):**
```json
{
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "first_name": "Test",
    "last_name": "User",
    "phone": "+254700000000",
    "role": "customer",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-10-02T20:30:00Z"
  }
}
```

#### 🔴 **Test 5: Authentication Error Cases** - [ ]

**5a. Invalid Credentials:**
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/auth/login`
- **Body:**
```json
{
  "email": "test@example.com",
  "password": "wrongpassword"
}
```
- **Expected Response (401):**
```json
{
  "error": "Invalid credentials"
}
```

**5b. Duplicate Registration:**
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/auth/register`
- **Body:** Same as Test 2
- **Expected Response (409):**
```json
{
  "error": "User with this email or username already exists"
}
```

**5c. Unauthorized Access:**
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/auth/profile`
- **Headers:** (No Authorization header)
- **Expected Response (401):**
```json
{
  "error": "Authorization header required"
}
```

---

### **Phase 3: Product Catalog Testing**

#### 🟢 **Test 6: Get All Products** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products`
- **Expected Response (200):**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Go Gopher T-Shirt",
      "slug": "go-gopher-tshirt",
      "description": "Official Go programming language mascot t-shirt",
      "category_id": 5,
      "base_price": 1500.00,
      "is_active": true,
      "is_featured": false
    }
  ],
  "total": 1
}
```

#### 🟢 **Test 7: Get Single Product** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products/1`
- **Expected Response (200):**
```json
{
  "id": 1,
  "name": "Go Gopher T-Shirt",
  "slug": "go-gopher-tshirt",
  "description": "Official Go programming language mascot t-shirt",
  "category_id": 5,
  "base_price": 1500.00,
  "is_active": true,
  "is_featured": false
}
```

#### 🟢 **Test 8: Get Categories** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/categories`
- **Expected Response (200):**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Clothing",
      "slug": "clothing",
      "description": "T-shirts, hoodies, and apparel",
      "parent_id": null,
      "is_active": true
    }
  ],
  "total": 1
}
```

#### 🔴 **Test 9: Product Not Found** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products/999`
- **Expected Response (404):**
```json
{
  "error": "Product not found"
}
```

---

### **Phase 4: Shopping Cart Testing**

**📝 Note:** For guest users, include the `X-Session-ID` header with a unique session identifier.

#### 🟢 **Test 10: Add Item to Cart (Guest User)** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/cart`
- **Headers:**
```
Content-Type: application/json
X-Session-ID: guest-session-12345
```
- **Body (JSON):**
```json
{
  "product_id": 1,
  "quantity": 2
}
```
- **Expected Response (200):**
```json
{
  "message": "Item added to cart successfully"
}
```

#### 🟢 **Test 11: Add Item to Cart (Authenticated User)** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/cart`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <your-jwt-token>
```
- **Body (JSON):**
```json
{
  "product_id": 1,
  "quantity": 3
}
```
- **Expected Response (200):**
```json
{
  "message": "Item added to cart successfully"
}
```

#### 🟢 **Test 12: Get Cart Contents (Guest)** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/cart`
- **Headers:**
```
X-Session-ID: guest-session-12345
```
- **Expected Response (200):**
```json
{
  "items": [
    {
      "id": 1,
      "session_id": "guest-session-12345",
      "product_id": 1,
      "quantity": 2,
      "product_name": "Go Gopher T-Shirt",
      "product_slug": "go-gopher-tshirt",
      "price": 1500.00
    }
  ],
  "total_items": 2,
  "subtotal": 3000.00
}
```

#### 🟢 **Test 13: Update Cart Item Quantity** - [ ]
- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/cart/1`
- **Headers:**
```
Content-Type: application/json
X-Session-ID: guest-session-12345
```
- **Body (JSON):**
```json
{
  "quantity": 5
}
```
- **Expected Response (200):**
```json
{
  "message": "Cart item updated successfully"
}
```

#### 🟢 **Test 14: Remove Item from Cart** - [ ]
- **Method:** `DELETE`
- **URL:** `http://localhost:8080/api/cart/1`
- **Headers:**
```
X-Session-ID: guest-session-12345
```
- **Expected Response (200):**
```json
{
  "message": "Item removed from cart successfully"
}
```

#### 🟢 **Test 15: Migrate Guest Cart (After Login)** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/cart/migrate`
- **Headers:**
```
Authorization: Bearer <your-jwt-token>
X-Session-ID: guest-session-12345
```
- **Expected Response (200):**
```json
{
  "message": "Guest cart migrated successfully"
}
```

---

### **Phase 5: Order Management Testing**

#### 🟢 **Test 16: Create Order from Cart** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/orders`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <your-jwt-token>
```
- **Body (JSON):**
```json
{
  "shipping_address": "123 Test Street, Nairobi, Kenya",
  "billing_address": "123 Test Street, Nairobi, Kenya",
  "notes": "Please deliver during business hours"
}
```
- **Expected Response (201):**
```json
{
  "message": "Order created successfully",
  "order": {
    "id": 1,
    "user_id": 1,
    "order_number": "ORD-1696273200",
    "status": "pending",
    "total_amount": 3000.00,
    "payment_status": "pending",
    "shipping_address": "123 Test Street, Nairobi, Kenya",
    "billing_address": "123 Test Street, Nairobi, Kenya",
    "notes": "Please deliver during business hours",
    "created_at": "2025-10-02T20:30:00Z",
    "updated_at": "2025-10-02T20:30:00Z",
    "items": [
      {
        "id": 1,
        "order_id": 1,
        "product_id": 1,
        "quantity": 2,
        "unit_price": 1500.00,
        "total_price": 3000.00,
        "product_name": "Go Gopher T-Shirt",
        "product_slug": "go-gopher-tshirt"
      }
    ]
  }
}
```

#### 🟢 **Test 17: Get Order by ID** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/orders/1`
- **Headers:**
```
Authorization: Bearer <your-jwt-token>
```
- **Expected Response (200):** Same as order object from Test 16

#### 🟢 **Test 18: Get User Orders** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/orders`
- **Headers:**
```
Authorization: Bearer <your-jwt-token>
```
- **Expected Response (200):**
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "order_number": "ORD-1696273200",
      "status": "pending",
      "total_amount": 3000.00,
      "payment_status": "pending",
      "created_at": "2025-10-02T20:30:00Z",
      "updated_at": "2025-10-02T20:30:00Z"
    }
  ],
  "total": 1
}
```

---

### **Phase 6: Admin Functionality Testing**

**⚠️ Important:** You need to manually promote a user to admin in the database first:

```sql
-- Use the correct schema-prefixed table name
UPDATE auth.users SET role = 'admin' WHERE email = 'test@example.com';
```

**📝 Note:** Since the migration to multi-schema architecture, the `users` table is now `auth.users`.

#### 🟢 **Test 19: Admin - Create Product** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/admin/products`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "name": "Python Developer Hoodie",
  "slug": "python-dev-hoodie",
  "description": "Comfortable hoodie for Python developers",
  "short_description": "Premium Python-themed hoodie",
  "category_id": 1,
  "base_price": 3500.00,
  "sku_prefix": "PY-HOOD",
  "is_featured": true,
  "weight": 0.5,
  "dimensions": "L x W x H"
}
```
- **Expected Response (201):**
```json
{
  "message": "Product created successfully",
  "product": {
    "id": 2,
    "name": "Python Developer Hoodie",
    "slug": "python-dev-hoodie",
    "description": "Comfortable hoodie for Python developers",
    "category_id": 1,
    "base_price": 3500.00,
    "is_active": true,
    "is_featured": true
  }
}
```

#### 🟢 **Test 20: Admin - Update Product** - [ ]
- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/admin/products/2`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "base_price": 3200.00,
  "is_featured": false
}
```
- **Expected Response (200):**
```json
{
  "message": "Product updated successfully",
  "product": {
    "id": 2,
    "name": "Python Developer Hoodie",
    "slug": "python-dev-hoodie",
    "description": "Comfortable hoodie for Python developers",
    "category_id": 1,
    "base_price": 3200.00,
    "is_active": true,
    "is_featured": false
  }
}
```

#### 🟢 **Test 21: Admin - Get All Products (Including Inactive)** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/admin/products`
- **Headers:**
```
Authorization: Bearer <admin-jwt-token>
```
- **Expected Response (200):**
```json
{
  "products": [
    {
      "id": 1,
      "name": "Go Gopher T-Shirt",
      "slug": "go-gopher-tshirt",
      "description": "Official Go programming language mascot t-shirt",
      "category_id": 5,
      "base_price": 1500.00,
      "is_active": true,
      "is_featured": false
    },
    {
      "id": 2,
      "name": "Python Developer Hoodie",
      "slug": "python-dev-hoodie",
      "description": "Comfortable hoodie for Python developers",
      "category_id": 1,
      "base_price": 3200.00,
      "is_active": true,
      "is_featured": false
    }
  ],
  "total": 2,
  "message": "All products (including inactive)"
}
```

#### 🟢 **Test 22: Admin - Delete Product (Soft Delete)** - [ ]
- **Method:** `DELETE`
- **URL:** `http://localhost:8080/api/admin/products/2`
- **Headers:**
```
Authorization: Bearer <admin-jwt-token>
```
- **Expected Response (200):**
```json
{
  "message": "Product deleted successfully"
}
```

#### 🟢 **Test 23: Admin - Create Category** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/admin/categories`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "name": "Clothing",
  "slug": "clothing",
  "description": "Apparel and clothing items",
  "sort_order": 1,
  "image_url": "https://example.com/clothing.jpg"
}
```
- **Expected Response (201):**
```json
{
  "message": "Category created successfully",
  "category": {
    "id": 1,
    "name": "Clothing",
    "slug": "clothing",
    "description": "Apparel and clothing items",
    "parent_id": null,
    "image_url": "https://example.com/clothing.jpg",
    "is_active": true,
    "sort_order": 1,
    "created_at": "2025-10-03T10:30:00Z",
    "updated_at": "2025-10-03T10:30:00Z"
  }
}
```

#### 🟢 **Test 24: Admin - Update Category** - [ ]
- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/admin/categories/1`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "name": "Premium Clothing",
  "description": "High-quality apparel and clothing items",
  "sort_order": 2
}
```
- **Expected Response (200):**
```json
{
  "message": "Category updated successfully",
  "category": {
    "id": 1,
    "name": "Premium Clothing",
    "slug": "clothing",
    "description": "High-quality apparel and clothing items",
    "parent_id": null,
    "image_url": "https://example.com/clothing.jpg",
    "is_active": true,
    "sort_order": 2,
    "created_at": "2025-10-03T10:30:00Z",
    "updated_at": "2025-10-03T10:32:00Z"
  }
}
```

#### 🟢 **Test 25: Admin - Get All Categories (Including Inactive)** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/admin/categories`
- **Headers:**
```
Authorization: Bearer <admin-jwt-token>
```
- **Expected Response (200):**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Premium Clothing",
      "slug": "clothing",
      "description": "High-quality apparel and clothing items",
      "parent_id": null,
      "image_url": "https://example.com/clothing.jpg",
      "is_active": true,
      "sort_order": 2,
      "created_at": "2025-10-03T10:30:00Z",
      "updated_at": "2025-10-03T10:32:00Z"
    }
  ],
  "total": 1,
  "message": "All categories (including inactive)"
}
```

#### 🟢 **Test 26: Admin - Delete Category (Soft Delete)** - [ ]
- **Method:** `DELETE`
- **URL:** `http://localhost:8080/api/admin/categories/1`
- **Headers:**
```
Authorization: Bearer <admin-jwt-token>
```
- **Expected Response (200):**
```json
{
  "message": "Category deleted successfully"
}
```

---

## 🖼️ **PRODUCT IMAGE MANAGEMENT**

#### 🟢 **Test 27: Create Product with Images** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/admin/products`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "name": "Test T-Shirt with Images",
  "slug": "test-t-shirt-with-images",
  "description": "A test t-shirt with multiple images for frontend integration",
  "short_description": "Test t-shirt with images",
  "category_id": 1,
  "base_price": 25.99,
  "sku_prefix": "TST",
  "is_featured": true,
  "weight": 0.2,
  "dimensions": "M",
  "images": [
    {
      "image_url": "https://example.com/images/test-tshirt-front.jpg",
      "image_path": "/uploads/test-tshirt-front.jpg",
      "image_type": "front",
      "alt_text": "Test T-Shirt Front View",
      "display_order": 1,
      "is_primary": true
    },
    {
      "image_url": "https://example.com/images/test-tshirt-back.jpg",
      "image_path": "/uploads/test-tshirt-back.jpg",
      "image_type": "back", 
      "alt_text": "Test T-Shirt Back View",
      "display_order": 2,
      "is_primary": false
    }
  ]
}
```
- **Expected Response (201):**
```json
{
  "message": "Product created successfully",
  "product": {
    "id": 1,
    "name": "Test T-Shirt with Images",
    "slug": "test-t-shirt-with-images",
    // ... other product fields
  },
  "images": [
    {
      "id": 1,
      "product_id": 1,
      "image_url": "https://example.com/images/test-tshirt-front.jpg",
      "is_primary": true,
      // ... other image fields
    },
    {
      "id": 2,
      "product_id": 1,
      "image_url": "https://example.com/images/test-tshirt-back.jpg",
      "is_primary": false,
      // ... other image fields
    }
  ],
  "images_created": 2
}
```

#### 🟢 **Test 28: Get Product Images (Public)** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products/1/images`
- **Headers:** None required
- **Expected Response (200):**
```json
{
  "images": [
    {
      "id": 1,
      "product_id": 1,
      "variant_id": null,
      "image_url": "https://example.com/images/test-tshirt-front.jpg",
      "image_path": "/uploads/test-tshirt-front.jpg",
      "image_type": "front",
      "alt_text": "Test T-Shirt Front View",
      "display_order": 1,
      "is_primary": true,
      "created_at": "2025-10-03T11:30:00.000Z"
    }
  ],
  "total": 1
}
```

#### 🟢 **Test 29: Admin - Add Image to Existing Product** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/admin/products/1/images`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "image_url": "https://example.com/images/test-tshirt-detail.jpg",
  "image_path": "/uploads/test-tshirt-detail.jpg",
  "image_type": "detail",
  "alt_text": "Test T-Shirt Detail View",
  "display_order": 3,
  "is_primary": false
}
```
- **Expected Response (201):**
```json
{
  "message": "Product image created successfully",
  "image": {
    "id": 3,
    "product_id": 1,
    "variant_id": null,
    "image_url": "https://example.com/images/test-tshirt-detail.jpg",
    "image_path": "/uploads/test-tshirt-detail.jpg",
    "image_type": "detail",
    "alt_text": "Test T-Shirt Detail View",
    "display_order": 3,
    "is_primary": false,
    "created_at": "2025-10-03T11:30:00.000Z"
  }
}
```

#### 🟢 **Test 30: Admin - Update Product Image** - [ ]
- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/admin/images/3`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "image_url": "https://example.com/images/test-tshirt-detail-updated.jpg",
  "image_path": "/uploads/test-tshirt-detail-updated.jpg",
  "image_type": "detail",
  "alt_text": "Test T-Shirt Detail View - Updated",
  "display_order": 3,
  "is_primary": false
}
```
- **Expected Response (200):**
```json
{
  "message": "Product image updated successfully",
  "image": {
    "id": 3,
    "product_id": 1,
    "image_url": "https://example.com/images/test-tshirt-detail-updated.jpg",
    // ... updated fields
  }
}
```

#### 🟢 **Test 31: Admin - Delete Product Image** - [ ]
- **Method:** `DELETE`
- **URL:** `http://localhost:8080/api/admin/images/3`
- **Headers:**
```
Authorization: Bearer <admin-jwt-token>
```
- **Expected Response (200):**
```json
{
  "message": "Product image deleted successfully"
}
```

---

## 📋 **ORDER MANAGEMENT**

#### 🟢 **Test 32: Admin - Get All Orders** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/admin/orders`
- **Headers:**
```
Authorization: Bearer <admin-jwt-token>
```
- **Expected Response (200):**
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "order_number": "ORD-1696273200",
      "status": "pending",
      "total_amount": 3000.00,
      "payment_status": "pending",
      "created_at": "2025-10-02T20:30:00Z",
      "updated_at": "2025-10-02T20:30:00Z"
    }
  ],
  "total": 1
}
```

#### 🟢 **Test 28: Admin - Update Order Status** - [ ]
- **Method:** `PUT`
- **URL:** `http://localhost:8080/api/admin/orders/1/status`
- **Headers:**
```
Content-Type: application/json
Authorization: Bearer <admin-jwt-token>
```
- **Body (JSON):**
```json
{
  "status": "processing",
  "payment_status": "paid",
  "payment_method": "mpesa"
}
```
- **Expected Response (200):**
```json
{
  "message": "Order status updated successfully"
}
```

#### 🔴 **Test 29: Admin Access Denied (Non-Admin User)** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/admin/products`
- **Headers:**
```
Authorization: Bearer <customer-jwt-token>
```
- **Expected Response (403):**
```json
{
  "error": "Admin access required"
}
```

---

### **Phase 7: Points System Testing**

#### 🟢 **Test 30: Get User Points** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/points`
- **Headers:**
```
Authorization: Bearer <your-jwt-token>
```
- **Expected Response (200):**
```json
{
  "user_id": 1,
  "points_balance": 0,
  "total_earned": 0,
  "total_spent": 0
}
```

---

## 🔧 **Error Testing Scenarios**

### **Validation Errors**

#### 🔴 **Test 27: Invalid Email Format** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/auth/register`
- **Body:**
```json
{
  "username": "testuser2",
  "email": "invalid-email",
  "password": "password123"
}
```

#### 🔴 **Test 28: Short Password** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/auth/register`
- **Body:**
```json
{
  "username": "testuser3",
  "email": "test3@example.com",
  "password": "123"
}
```

#### 🔴 **Test 29: Missing Required Fields** - [ ]
- **Method:** `POST`
- **URL:** `http://localhost:8080/api/cart`
- **Body:**
```json
{
  "quantity": 2
}
```

#### 🔴 **Test 30: Invalid Product ID** - [ ]
- **Method:** `GET`
- **URL:** `http://localhost:8080/api/products/abc`

---

## 📊 **Postman Collection Setup**

### **Environment Variables**
Create a Postman environment with these variables:

```
base_url: http://localhost:8080
user_token: (set after login)
admin_token: (set after admin login)
session_id: guest-session-12345
```

### **Collection Structure**

```
📁 Merch Ke API Tests
├── 📁 01. Health Check
│   └── Health Check
├── 📁 02. Authentication
│   ├── Register User
│   ├── Login User
│   ├── Get Profile
│   └── Auth Error Cases
├── 📁 03. Product Catalog
│   ├── Get All Products
│   ├── Get Single Product
│   ├── Get Categories
│   └── Product Not Found
├── 📁 04. Shopping Cart
│   ├── Add to Cart (Guest)
│   ├── Add to Cart (User)
│   ├── Get Cart
│   ├── Update Cart
│   ├── Remove from Cart
│   └── Migrate Cart
├── 📁 05. Orders
│   ├── Create Order
│   ├── Get Order
│   └── Get User Orders
├── 📁 06. Admin Functions
│   ├── Create Product
│   ├── Update Product
│   ├── Delete Product
│   ├── Get All Products
│   ├── Create Category
│   ├── Update Category
│   ├── Delete Category
│   ├── Get All Categories
│   ├── Get All Orders
│   ├── Update Order Status
│   └── Access Denied
└── 📁 07. Points System
    └── Get User Points
```

---

## 🎯 **Testing Checklist**

### **Functional Testing**
- [ ] User registration and login
- [ ] JWT token generation and validation
- [ ] Product catalog browsing
- [ ] Category management (public endpoints)
- [ ] Cart operations (add, update, remove)
- [ ] Guest vs authenticated user flows
- [ ] Order creation and management
- [ ] Admin product management
- [ ] Admin category management
- [ ] **Admin product image management** 🆕
- [ ] **Product creation with images** 🆕
- [ ] Admin order management
- [ ] Points system functionality

### **Security Testing**
- [ ] Unauthorized access attempts
- [ ] Invalid token handling
- [ ] Admin privilege enforcement
- [ ] SQL injection protection
- [ ] Input validation

### **Error Handling**
- [ ] Invalid request formats
- [ ] Missing required fields
- [ ] Non-existent resources
- [ ] Database connection errors
- [ ] Validation errors

### **Performance Testing**
- [ ] Response times under normal load
- [ ] Database query efficiency
- [ ] Memory usage patterns

---

## 🚨 **Common Issues & Solutions**

### **Database Connection Issues**
```bash
# Check if PostgreSQL is running
systemctl status postgresql

# Test connection manually
psql -U merch-ke-admin -d merch-ke-db -h localhost

# View schema structure
psql -U merch-ke-admin -d merch-ke-db -c "\dn"  # List schemas
psql -U merch-ke-admin -d merch-ke-db -c "\dt auth.*"  # List auth tables
psql -U merch-ke-admin -d merch-ke-db -c "\dt catalog.*"  # List catalog tables
psql -U merch-ke-admin -d merch-ke-db -c "\dt orders.*"  # List orders tables
```

### **Multi-Schema Database Architecture**
The API uses a modern multi-schema PostgreSQL design:

**Schema Organization:**
- **`auth.users`** - User accounts (not just `users`)
- **`auth.user_addresses`** - User shipping addresses
- **`auth.user_points`** - Loyalty points system
- **`catalog.products`** - Product catalog
- **`catalog.categories`** - Product categories
- **`orders.orders`** - Customer orders
- **`orders.cart_items`** - Shopping cart

**Important SQL Commands:**
```sql
-- Correct way to query users
SELECT * FROM auth.users WHERE email = 'user@example.com';

-- Correct way to update user role
UPDATE auth.users SET role = 'admin' WHERE email = 'user@example.com';

-- View all products
SELECT * FROM catalog.products WHERE is_active = true;
```

### **JWT Token Issues**
- Ensure JWT_SECRET is set in environment
- Check token expiration (24 hours default)
- Verify Bearer token format in Authorization header

### **CORS Issues**
The API includes CORS middleware, but if you encounter issues:
- Ensure requests are from allowed origins
- Check browser console for CORS errors

### **Session ID for Guest Users**
- Always include `X-Session-ID` header for guest operations
- Use consistent session ID throughout guest shopping flow

---

## 📝 **Notes for Production**

1. **Environment Variables**: Update all sensitive data in `.env`
2. **Database**: Use proper PostgreSQL instance with backup strategy
3. **JWT Secret**: Use cryptographically secure random string
4. **Rate Limiting**: Implement rate limiting for production
5. **Logging**: Configure proper logging and monitoring
6. **HTTPS**: Always use HTTPS in production
7. **Input Validation**: Additional validation layers recommended
8. **Error Messages**: Don't expose sensitive information in error messages
9. **Schema Privileges**: Set proper schema-level permissions for security
10. **Database Migrations**: Use proper migration tools for schema updates

---

## 📊 **Testing Summary**

### **Test Count Overview:**
- **Total Tests**: 30+ test cases
- **Success Tests** (🟢): ~26 tests  
- **Error Tests** (🔴): ~4 tests

### **Quick Test Categories:**
- **Phase 1**: Health Check (1 test)
- **Phase 2**: Authentication (4 tests)  
- **Phase 3**: Product Catalog (3 tests)
- **Phase 4**: Shopping Cart (6 tests)
- **Phase 5**: Orders (3 tests)
- **Phase 6**: Admin Functions (13 tests)
- **Phase 7**: Points System (1 test)

### **Progress Calculation:**
**Completion Rate**: [Passed Tests] / [Total Tests] × 100%

Example: If you complete 25 out of 30 tests = 83% completion

---

## 🎉 **Success Criteria**

Your API is ready for production when:
- ✅ All test cases pass
- ✅ No temporary/development code remains
- ✅ Proper error handling for all edge cases
- ✅ Security measures implemented
- ✅ Performance benchmarks met
- ✅ Multi-schema database architecture is properly implemented
- ✅ Schema-level permissions are configured
- ✅ Documentation is complete and updated for new schema structure

Happy Testing! 🚀