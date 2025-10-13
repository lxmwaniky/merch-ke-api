# Merch Ke API

A production-ready Go REST API for an e-commerce merch store, featuring authentication, product catalog, shopping cart, orders, and loyalty points. Built with Go 1.25 and Fiber v2, backed by a multi-schema PostgreSQL database.

## 🚀 Features

### Core Functionality
- **Authentication & Authorization** - JWT-based auth with role-based access control (customer/admin)
- **Product Catalog** - Multi-category product management with variants and images
- **Shopping Cart** - Session-aware cart for both guests and authenticated users
- **Order Management** - Complete order lifecycle with status tracking
- **Loyalty Points** - Points accumulation and transaction history
- **Admin Dashboard** - Full CRUD operations for products, categories, and orders

### Technical Highlights
- Multi-schema PostgreSQL architecture (`auth`, `catalog`, `orders`)
- Cloud-native design (Google Cloud Run + Cloud SQL)
- Guest-to-user cart migration on login
- Secure password hashing with bcrypt
- RESTful API design with comprehensive error handling

## 🛠️ Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.25 |
| **Web Framework** | Fiber v2 |
| **Database** | PostgreSQL 15 |
| **Authentication** | JWT (HS256) |
| **Password Hashing** | bcrypt |
| **Environment Config** | godotenv |
| **Deployment** | Docker + Google Cloud Run |
| **Database Hosting** | Google Cloud SQL |

## 📁 Project Structure

```
merch-ke-api/
├── main.go                 # Application entry point and route definitions
├── auth.go                 # Authentication handlers and JWT middleware
├── handlers.go             # Business logic for catalog, cart, and orders
├── models.go               # Data models and database helpers
├── database.go             # Database connection and configuration
├── Dockerfile              # Multi-stage build for production
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── .env.example            # Environment variables template
├── .gitignore              # Git ignore rules
├── TESTING.md              # Testing documentation and guidelines
├── *_test.go               # Unit tests for auth, models, and handlers
└── database/
    └── schema.sql          # Complete database schema with multi-schema design
```

## 🚦 Getting Started

### Prerequisites

- **Go 1.25+** - [Download](https://go.dev/dl/)
- **PostgreSQL 15+** - Local installation, Docker, or cloud instance
- **psql** - PostgreSQL CLI for running migrations
- **Git** - For cloning the repository

### 1. Clone the Repository

```bash
git clone https://github.com/lxmwaniky/merch-ke-api.git
cd merch-ke-api
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Environment Variables

Create a `.env` file from the template:

```bash
cp .env.example .env
```

Configure the following variables in `.env`:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=merch_ke_db
DB_SSLMODE=disable

# JWT Secret (generate a strong random string)
JWT_SECRET=your-super-secret-jwt-key-min-32-chars

# Server Port (optional, defaults to 8080)
PORT=8080
```

**Security Note:** Use a strong, random JWT secret in production. Generate one using:
```bash
openssl rand -hex 32
```

### 4. Create and Initialize Database

```bash
# Create the database
psql -U postgres -c "CREATE DATABASE merch_ke_db"

# Run the schema migration
psql -U postgres -d merch_ke_db -f database/schema.sql
```

This creates three schemas:
- `auth` - Users, addresses, and loyalty points
- `catalog` - Products, categories, variants, and images
- `orders` - Shopping carts, orders, and order items

### 5. Run the Application

#### Standard Mode
```bash
go run .
```

#### With Hot Reload (Development)

Install Air for automatic reloading:
```bash
go install github.com/cosmtrek/air@latest
air
```

The API will start on `http://localhost:8080` (or your configured PORT).

### 6. Verify Installation

Test the health endpoint:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "Merch Ke API"
}
```

### 7. Run Tests

Run the unit tests to verify everything is working:

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

All tests should pass ✅ See [TESTING.md](TESTING.md) for detailed testing documentation.

---

## 🧪 Testing

This project includes comprehensive unit tests covering:
- Authentication & JWT token generation
- Password hashing and verification
- Data models and validation
- HTTP handlers and middleware
- Input validation and error handling

**Current Coverage**: 3.8% (unit tests for core functionality)

Run tests:
```bash
go test -v ./...
  "service": "Merch Ke API"
}
```

## 🧪 Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 📊 Database Schema

The application uses a multi-schema PostgreSQL design for logical separation:

### `auth` Schema
- `auth.users` - User accounts and authentication
- `auth.user_addresses` - Shipping/billing addresses
- `auth.user_points` - Current loyalty points balance
- `auth.points_transactions` - Points transaction history

### `catalog` Schema
- `catalog.categories` - Product categories (hierarchical)
- `catalog.products` - Product catalog
- `catalog.product_variants` - Size, color, SKU variations
- `catalog.product_images` - Product image URLs

### `orders` Schema
- `orders.cart_items` - Authenticated user shopping carts
- `orders.guest_cart_items` - Guest session shopping carts
- `orders.orders` - Order records
- `orders.order_items` - Line items in orders

All tables include appropriate indexes, foreign keys, and constraints for data integrity.

## 🔌 API Overview

The API provides RESTful endpoints organized by functionality.

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check |
| `POST` | `/api/auth/register` | Create new user account |
| `POST` | `/api/auth/login` | Authenticate and get JWT token |
| `GET` | `/api/products` | List all products |
| `GET` | `/api/products/:id` | Get single product details |
| `GET` | `/api/categories` | List all categories |

### Protected Endpoints (Requires JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/auth/profile` | Get current user profile |
| `POST` | `/api/cart` | Add item to cart |
| `GET` | `/api/cart` | Get cart contents |
| `PUT` | `/api/cart/:productId` | Update cart item quantity |
| `DELETE` | `/api/cart/:productId` | Remove item from cart |
| `POST` | `/api/orders` | Create order from cart |
| `GET` | `/api/orders/:id` | Get order details |
| `GET` | `/api/orders` | List user's orders |

### Admin Endpoints (Requires Admin Role)

All admin endpoints are prefixed with `/api/admin` and require authentication with admin privileges.

- Product management (CRUD)
- Category management (CRUD)
- Product image management
- Order status updates
- View all orders

## 🔐 Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. **Register** or **Login** to receive a JWT token
2. Include the token in subsequent requests:
   ```
   Authorization: Bearer <your-jwt-token>
   ```
3. Tokens contain user ID and role (customer/admin)
4. Tokens expire after a configured duration

### Guest vs Authenticated Carts

- **Guest users**: Use `X-Session-ID` header with a unique session identifier
- **Authenticated users**: Carts are automatically tied to user account
- **Cart migration**: When a guest logs in, their cart merges with their account cart

## 🐛 Troubleshooting

### Database Connection Issues

**Problem:** `Failed to connect to database`

**Solutions:**
- Verify PostgreSQL is running: `pg_isready`
- Check credentials in `.env` file
- Ensure database exists: `psql -U postgres -l`
- For Cloud SQL, verify connection name and service account permissions

### JWT Token Errors

**Problem:** `Authorization header required` or `Invalid token`

**Solutions:**
- Ensure `Authorization: Bearer <token>` header is included
- Verify JWT_SECRET matches between token generation and validation
- Check token hasn't expired
- Re-login to get a fresh token

### Port Already in Use

**Problem:** `address already in use`

**Solutions:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change PORT in .env file
```

### Schema Migration Fails

**Problem:** Errors when running `schema.sql`

**Solutions:**
- Drop and recreate database: `DROP DATABASE merch_ke_db; CREATE DATABASE merch_ke_db;`
- Check PostgreSQL version compatibility (requires 15+)
- Ensure you have proper permissions

## 📝 Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_HOST` | Yes | - | Database hostname or socket path |
| `DB_PORT` | Yes | `5432` | Database port |
| `DB_USER` | Yes | - | Database username |
| `DB_PASSWORD` | Yes | - | Database password |
| `DB_NAME` | Yes | - | Database name |
| `DB_SSLMODE` | No | `disable` | SSL mode (`disable`, `require`, `verify-full`) |
| `JWT_SECRET` | Yes | - | Secret key for JWT signing (min 32 chars) |
| `PORT` | No | `8080` | HTTP server port |

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

For questions or issues:
- Contact: lekko254@gmail.com

---

**Status**: Active Development 🚧  
**Version**: 1.0.0  
**Last Updated**: October 2025