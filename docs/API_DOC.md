# Merch Ke API Documentation

> **For Frontend Teams**: Complete API reference with request/response examples for building the Merch Ke frontend application.

## 🌐 Base URL

**Production API**: `https://merch-ke-api-779644650318.us-central1.run.app`

All endpoints should be prefixed with this base URL.

Example:
```
GET https://merch-ke-api-779644650318.us-central1.run.app/health
```

## 🔑 Authentication

Most endpoints require authentication via JWT tokens. After registering or logging in, include the token in your requests:

```
Authorization: Bearer <your-jwt-token>
```

### Guest Users

For guest (non-authenticated) users accessing the shopping cart, include a session identifier:

```
X-Session-ID: <unique-session-id>
```

Generate this on the frontend (e.g., UUID) and persist it in localStorage/cookies until the user registers or logs in.

## 📊 Database Schema Overview

The API uses a multi-schema PostgreSQL architecture:

| Schema | Tables | Purpose |
|--------|--------|---------|
| `auth` | users, user_addresses, user_points, points_transactions | Authentication and user management |
| `catalog` | products, categories, product_variants, product_images | Product catalog and categorization |
| `orders` | cart_items, guest_cart_items, orders, order_items | Shopping cart and order management |

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

## 📑 Table of Contents

1. [Health Check](#health-check)
2. [Authentication](#authentication-endpoints)
   - Register
   - Login
   - Get Profile
3. [Product Catalog](#product-catalog-endpoints)
   - List Products
   - Get Single Product
   - List Categories
   - Get Product Images
4. [Shopping Cart](#shopping-cart-endpoints)
   - Add to Cart
   - Get Cart
   - Update Cart Item
   - Remove from Cart
   - Migrate Guest Cart
5. [Orders](#order-endpoints)
   - Create Order
   - Get Order Details
   - List User Orders
6. [Loyalty Points](#loyalty-points-endpoints)
   - Get User Points
7. [Admin Endpoints](#admin-endpoints)
   - Product Management
   - Category Management
   - Image Management
   - Order Management

---

## Health Check

### GET /health

Check if the API is running and healthy.

**Request:**
```http
GET /health
```

**Response:** `200 OK`
```json
{
  "status": "healthy",
  "service": "Merch Ke API"
}
```

---

## Authentication Endpoints

### POST /api/auth/register

Create a new user account.

**Request:**
```http
POST /api/auth/register
Content-Type: application/json
```

**Body:**
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+254712345678"
}
```

**Response:** `201 Created`
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+254712345678",
    "role": "customer",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-10-13T09:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `409 Conflict` - Email or username already exists
- `400 Bad Request` - Invalid input data

---

### POST /api/auth/login

Authenticate and receive a JWT token.

**Request:**
```http
POST /api/auth/login
Content-Type: application/json
```

**Body:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:** `200 OK`
```json
{
  "message": "Login successful",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+254712345678",
    "role": "customer",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-10-13T09:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Errors:**
- `401 Unauthorized` - Invalid credentials
- `400 Bad Request` - Missing email or password

---

### GET /api/auth/profile

Get the authenticated user's profile information.

**Request:**
```http
GET /api/auth/profile
Authorization: Bearer <jwt-token>
```

**Response:** `200 OK`
```json
{
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+254712345678",
    "role": "customer",
    "is_active": true,
    "email_verified": false,
    "created_at": "2025-10-13T09:00:00Z"
  }
}
```

**Errors:**
- `401 Unauthorized` - Missing or invalid token

---

## Product Catalog Endpoints

### GET /api/products

Retrieve all active products.

**Request:**
```http
GET /api/products
```

**Response:** `200 OK`
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
      "is_featured": false,
      "created_at": "2025-10-10T10:00:00Z"
    },
    {
      "id": 2,
      "name": "Docker Whale Hoodie",
      "slug": "docker-whale-hoodie",
      "description": "Comfortable hoodie with Docker logo",
      "category_id": 5,
      "base_price": 3500.00,
      "is_active": true,
      "is_featured": true,
      "created_at": "2025-10-11T14:30:00Z"
    }
  ],
  "total": 2
}
```

---

### GET /api/products/:id

Get detailed information about a single product.

**Request:**
```http
GET /api/products/1
```

**Response:** `200 OK`
```json
{
  "id": 1,
  "name": "Go Gopher T-Shirt",
  "slug": "go-gopher-tshirt",
  "description": "Official Go programming language mascot t-shirt. Available in multiple sizes and colors.",
  "category_id": 5,
  "base_price": 1500.00,
  "is_active": true,
  "is_featured": false,
  "created_at": "2025-10-10T10:00:00Z",
  "updated_at": "2025-10-10T10:00:00Z"
}
```

**Errors:**
- `404 Not Found` - Product doesn't exist
- `400 Bad Request` - Invalid product ID

---

### GET /api/categories

List all product categories.

**Request:**
```http
GET /api/categories
```

**Response:** `200 OK`
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Clothing",
      "slug": "clothing",
      "description": "T-shirts, hoodies, and apparel",
      "parent_id": null,
      "is_active": true,
      "created_at": "2025-10-01T00:00:00Z"
    },
    {
      "id": 5,
      "name": "Tech Apparel",
      "slug": "tech-apparel",
      "description": "Programming and tech-themed clothing",
      "parent_id": 1,
      "is_active": true,
      "created_at": "2025-10-02T00:00:00Z"
    }
  ],
  "total": 2
}
```

---

### GET /api/products/:productId/images

Get all images for a specific product.

**Request:**
```http
GET /api/products/1/images
```

**Response:** `200 OK`
```json
{
  "images": [
    {
      "id": 1,
      "product_id": 1,
      "image_url": "https://example.com/images/gopher-tshirt-front.jpg",
      "alt_text": "Go Gopher T-Shirt Front View",
      "display_order": 1,
      "is_primary": true
    },
    {
      "id": 2,
      "product_id": 1,
      "image_url": "https://example.com/images/gopher-tshirt-back.jpg",
      "alt_text": "Go Gopher T-Shirt Back View",
      "display_order": 2,
      "is_primary": false
    }
  ]
}
```

---

## Shopping Cart Endpoints

### POST /api/cart

Add an item to the shopping cart. Works for both authenticated and guest users.

**Request (Authenticated User):**
```http
POST /api/cart
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Request (Guest User):**
```http
POST /api/cart
X-Session-ID: <unique-session-id>
Content-Type: application/json
```

**Body:**
```json
{
  "product_id": 1,
  "quantity": 2
}
```

**Response:** `200 OK`
```json
{
  "message": "Item added to cart successfully"
}
```

**Errors:**
- `400 Bad Request` - Invalid product_id or quantity
- `404 Not Found` - Product doesn't exist
- `401 Unauthorized` - Missing both JWT and session ID

---

### GET /api/cart

Retrieve the current user's or guest's shopping cart.

**Request (Authenticated):**
```http
GET /api/cart
Authorization: Bearer <jwt-token>
```

**Request (Guest):**
```http
GET /api/cart
X-Session-ID: <unique-session-id>
```

**Response:** `200 OK`
```json
{
  "items": [
    {
      "id": 1,
      "product_id": 1,
      "quantity": 2,
      "product_name": "Go Gopher T-Shirt",
      "product_slug": "go-gopher-tshirt",
      "price": 1500.00,
      "subtotal": 3000.00
    },
    {
      "id": 2,
      "product_id": 2,
      "quantity": 1,
      "product_name": "Docker Whale Hoodie",
      "product_slug": "docker-whale-hoodie",
      "price": 3500.00,
      "subtotal": 3500.00
    }
  ],
  "total_items": 3,
  "subtotal": 6500.00
}
```

---

### PUT /api/cart/:productId

Update the quantity of an item in the cart.

**Request:**
```http
PUT /api/cart/1
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "quantity": 5
}
```

**Response:** `200 OK`
```json
{
  "message": "Cart item updated successfully"
}
```

**Errors:**
- `400 Bad Request` - Invalid quantity (must be > 0)
- `404 Not Found` - Item not in cart

---

### DELETE /api/cart/:productId

Remove an item from the cart.

**Request:**
```http
DELETE /api/cart/1
Authorization: Bearer <jwt-token>
```

**Response:** `200 OK`
```json
{
  "message": "Item removed from cart successfully"
}
```

---

### POST /api/cart/migrate

Migrate guest cart items to authenticated user's cart after login/registration.

**Request:**
```http
POST /api/cart/migrate
Authorization: Bearer <jwt-token>
X-Session-ID: <guest-session-id>
```

**Response:** `200 OK`
```json
{
  "message": "Cart migrated successfully",
  "items_migrated": 3
}
```

---

## Order Endpoints

### POST /api/orders

Create a new order from the current cart items.

**Request:**
```http
POST /api/orders
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "shipping_address_id": 1,
  "payment_method": "mpesa",
  "notes": "Please deliver between 9 AM - 5 PM"
}
```

**Response:** `201 Created`
```json
{
  "message": "Order created successfully",
  "order": {
    "id": 123,
    "order_number": "ORD-20251013-0123",
    "user_id": 1,
    "total_amount": 6500.00,
    "status": "pending",
    "payment_method": "mpesa",
    "created_at": "2025-10-13T10:30:00Z"
  }
}
```

**Errors:**
- `400 Bad Request` - Empty cart or invalid address
- `401 Unauthorized` - Not authenticated (guest users cannot place orders)

---

### GET /api/orders/:id

Get details of a specific order.

**Request:**
```http
GET /api/orders/123
Authorization: Bearer <jwt-token>
```

**Response:** `200 OK`
```json
{
  "id": 123,
  "order_number": "ORD-20251013-0123",
  "user_id": 1,
  "total_amount": 6500.00,
  "status": "pending",
  "payment_method": "mpesa",
  "notes": "Please deliver between 9 AM - 5 PM",
  "items": [
    {
      "id": 1,
      "product_id": 1,
      "product_name": "Go Gopher T-Shirt",
      "quantity": 2,
      "price": 1500.00,
      "subtotal": 3000.00
    },
    {
      "id": 2,
      "product_id": 2,
      "product_name": "Docker Whale Hoodie",
      "quantity": 1,
      "price": 3500.00,
      "subtotal": 3500.00
    }
  ],
  "created_at": "2025-10-13T10:30:00Z",
  "updated_at": "2025-10-13T10:30:00Z"
}
```

**Errors:**
- `404 Not Found` - Order doesn't exist
- `403 Forbidden` - Order belongs to another user

---

### GET /api/orders

List all orders for the authenticated user.

**Request:**
```http
GET /api/orders
Authorization: Bearer <jwt-token>
```

**Response:** `200 OK`
```json
{
  "orders": [
    {
      "id": 123,
      "order_number": "ORD-20251013-0123",
      "total_amount": 6500.00,
      "status": "pending",
      "created_at": "2025-10-13T10:30:00Z"
    },
    {
      "id": 122,
      "order_number": "ORD-20251012-0122",
      "total_amount": 4200.00,
      "status": "delivered",
      "created_at": "2025-10-12T14:20:00Z"
    }
  ],
  "total": 2
}
```

---

## Loyalty Points Endpoints

### GET /api/points

Get the authenticated user's loyalty points balance and recent transactions.

**Request:**
```http
GET /api/points
Authorization: Bearer <jwt-token>
```

**Response:** `200 OK`
```json
{
  "balance": 150,
  "transactions": [
    {
      "id": 1,
      "points": 100,
      "transaction_type": "earned",
      "description": "Order #ORD-20251013-0123",
      "created_at": "2025-10-13T10:35:00Z"
    },
    {
      "id": 2,
      "points": 50,
      "transaction_type": "earned",
      "description": "Welcome bonus",
      "created_at": "2025-10-13T09:00:00Z"
    }
  ]
}
```

---

## Admin Endpoints

All admin endpoints require authentication with an admin role JWT token.

### Product Management

#### POST /api/admin/products

Create a new product.

**Request:**
```http
POST /api/admin/products
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Kubernetes Mug",
  "slug": "kubernetes-mug",
  "description": "Ceramic coffee mug with Kubernetes logo",
  "category_id": 3,
  "base_price": 800.00,
  "is_active": true,
  "is_featured": false
}
```

**Response:** `201 Created`
```json
{
  "message": "Product created successfully",
  "product": {
    "id": 10,
    "name": "Kubernetes Mug",
    "slug": "kubernetes-mug",
    "description": "Ceramic coffee mug with Kubernetes logo",
    "category_id": 3,
    "base_price": 800.00,
    "is_active": true,
    "is_featured": false,
    "created_at": "2025-10-13T11:00:00Z"
  }
}
```

---

#### PUT /api/admin/products/:id

Update an existing product.

**Request:**
```http
PUT /api/admin/products/10
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Kubernetes Coffee Mug",
  "base_price": 900.00,
  "is_featured": true
}
```

**Response:** `200 OK`
```json
{
  "message": "Product updated successfully"
}
```

---

#### DELETE /api/admin/products/:id

Delete a product (soft delete).

**Request:**
```http
DELETE /api/admin/products/10
Authorization: Bearer <admin-jwt-token>
```

**Response:** `200 OK`
```json
{
  "message": "Product deleted successfully"
}
```

---

#### GET /api/admin/products

List all products (including inactive ones) for admin management.

**Request:**
```http
GET /api/admin/products
Authorization: Bearer <admin-jwt-token>
```

**Response:** `200 OK`
```json
{
  "products": [
    {
      "id": 1,
      "name": "Go Gopher T-Shirt",
      "slug": "go-gopher-tshirt",
      "category_id": 5,
      "base_price": 1500.00,
      "is_active": true,
      "is_featured": false
    }
  ],
  "total": 1
}
```

---

### Category Management

#### POST /api/admin/categories

Create a new category.

**Request:**
```http
POST /api/admin/categories
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "name": "Accessories",
  "slug": "accessories",
  "description": "Bags, stickers, and other accessories",
  "parent_id": null,
  "is_active": true
}
```

**Response:** `201 Created`
```json
{
  "message": "Category created successfully",
  "category": {
    "id": 10,
    "name": "Accessories",
    "slug": "accessories",
    "description": "Bags, stickers, and other accessories",
    "parent_id": null,
    "is_active": true,
    "created_at": "2025-10-13T11:30:00Z"
  }
}
```

---

#### PUT /api/admin/categories/:id

Update a category.

**Request:**
```http
PUT /api/admin/categories/10
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "description": "Bags, stickers, mugs, and other merchandise"
}
```

**Response:** `200 OK`
```json
{
  "message": "Category updated successfully"
}
```

---

#### DELETE /api/admin/categories/:id

Delete a category.

**Request:**
```http
DELETE /api/admin/categories/10
Authorization: Bearer <admin-jwt-token>
```

**Response:** `200 OK`
```json
{
  "message": "Category deleted successfully"
}
```

---

### Product Image Management

#### POST /api/admin/products/:productId/images

Add an image to a product.

**Request:**
```http
POST /api/admin/products/1/images
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "image_url": "https://example.com/images/product-123.jpg",
  "alt_text": "Product front view",
  "display_order": 1,
  "is_primary": true
}
```

**Response:** `201 Created`
```json
{
  "message": "Image added successfully",
  "image": {
    "id": 15,
    "product_id": 1,
    "image_url": "https://example.com/images/product-123.jpg",
    "alt_text": "Product front view",
    "display_order": 1,
    "is_primary": true
  }
}
```

---

#### PUT /api/admin/images/:imageId

Update a product image.

**Request:**
```http
PUT /api/admin/images/15
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "alt_text": "Product front view - updated",
  "display_order": 2
}
```

**Response:** `200 OK`
```json
{
  "message": "Image updated successfully"
}
```

---

#### DELETE /api/admin/images/:imageId

Delete a product image.

**Request:**
```http
DELETE /api/admin/images/15
Authorization: Bearer <admin-jwt-token>
```

**Response:** `200 OK`
```json
{
  "message": "Image deleted successfully"
}
```

---

### Order Management

#### GET /api/admin/orders

View all orders in the system.

**Request:**
```http
GET /api/admin/orders
Authorization: Bearer <admin-jwt-token>
```

**Response:** `200 OK`
```json
{
  "orders": [
    {
      "id": 123,
      "order_number": "ORD-20251013-0123",
      "user_id": 1,
      "user_email": "john@example.com",
      "total_amount": 6500.00,
      "status": "pending",
      "created_at": "2025-10-13T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

#### PUT /api/admin/orders/:id/status

Update order status.

**Request:**
```http
PUT /api/admin/orders/123/status
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json
```

**Body:**
```json
{
  "status": "processing"
}
```

**Valid statuses:** `pending`, `processing`, `shipped`, `delivered`, `cancelled`

**Response:** `200 OK`
```json
{
  "message": "Order status updated successfully"
}
```

---

## Error Responses

All endpoints return consistent error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request data",
  "details": "Quantity must be greater than 0"
}
```

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 403 Forbidden
```json
{
  "error": "Access denied. Admin privileges required."
}
```

### 404 Not Found
```json
{
  "error": "Product not found"
}
```

### 409 Conflict
```json
{
  "error": "User with this email already exists"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "details": "An unexpected error occurred"
}
```

---

## Rate Limiting

The API implements rate limiting to prevent abuse. Current limits:

- **Authenticated requests**: 100 requests per minute
- **Unauthenticated requests**: 20 requests per minute

When rate limited, you'll receive a `429 Too Many Requests` response:
```json
{
  "error": "Rate limit exceeded. Please try again later."
}
```

---

## CORS

The API supports Cross-Origin Resource Sharing (CORS) for frontend applications. All origins are currently allowed in development. In production, only whitelisted domains will be permitted.

---

## Support & Contact

For API support or to report issues:
- **GitHub Issues**: [https://github.com/lxmwaniky/merch-ke-api/issues](https://github.com/lxmwaniky/merch-ke-api/issues)
- **Email**: lekko254@gmail.com

---

**Last Updated**: October 13, 2025  
**API Version**: 1.0.0
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