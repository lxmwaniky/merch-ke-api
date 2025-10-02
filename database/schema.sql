-- Merch Ke - Multi-Schema Database Design
-- Organized schemas: catalog, auth, orders

-- =====================================================
-- Create Schemas
-- =====================================================
CREATE SCHEMA IF NOT EXISTS catalog;
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS orders;

-- =====================================================
-- CATALOG SCHEMA - Products, Categories, Variants
-- =====================================================

CREATE TABLE catalog.categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES catalog.categories(id),
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE catalog.products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    short_description VARCHAR(500),
    category_id INTEGER REFERENCES catalog.categories(id),
    base_price DECIMAL(10,2) NOT NULL,
    sku_prefix VARCHAR(20),
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    weight DECIMAL(8,2),
    dimensions VARCHAR(100),
    meta_title VARCHAR(255),
    meta_description VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE catalog.product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES catalog.products(id) ON DELETE CASCADE,
    sku VARCHAR(100) UNIQUE NOT NULL,
    size VARCHAR(20),
    color VARCHAR(50),
    material VARCHAR(100),
    price_adjustment DECIMAL(10,2) DEFAULT 0.00,
    stock_quantity INTEGER DEFAULT 0,
    low_stock_threshold INTEGER DEFAULT 5,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE catalog.product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES catalog.products(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES catalog.product_variants(id) ON DELETE SET NULL,
    image_url VARCHAR(500) NOT NULL,
    image_path VARCHAR(500) NOT NULL,
    image_type VARCHAR(50) DEFAULT 'gallery',
    alt_text VARCHAR(255),
    display_order INTEGER DEFAULT 1,
    file_size INTEGER,
    width INTEGER,
    height INTEGER,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- AUTH SCHEMA - Users, Addresses, Points
-- =====================================================

CREATE TABLE auth.users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    role VARCHAR(20) DEFAULT 'customer',
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auth.user_addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES auth.users(id) ON DELETE CASCADE,
    type VARCHAR(20) DEFAULT 'shipping',
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    company VARCHAR(100),
    address_line_1 VARCHAR(255) NOT NULL,
    address_line_2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    county VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(100) DEFAULT 'Kenya',
    phone VARCHAR(20),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auth.user_points (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES auth.users(id) ON DELETE CASCADE,
    points_balance INTEGER DEFAULT 0,
    total_earned INTEGER DEFAULT 0,
    total_spent INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auth.points_transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES auth.users(id) ON DELETE CASCADE,
    order_id INTEGER REFERENCES orders.orders(id),
    transaction_type VARCHAR(20) NOT NULL,
    points INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- ORDERS SCHEMA - Orders, Carts, Order Items
-- =====================================================

CREATE TABLE orders.orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES auth.users(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    payment_status VARCHAR(20) DEFAULT 'pending',
    payment_method VARCHAR(50),
    payment_reference VARCHAR(100),
    subtotal DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    shipping_amount DECIMAL(10,2) DEFAULT 0.00,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL,
    shipping_first_name VARCHAR(100),
    shipping_last_name VARCHAR(100),
    shipping_company VARCHAR(100),
    shipping_address_line_1 VARCHAR(255),
    shipping_address_line_2 VARCHAR(255),
    shipping_city VARCHAR(100),
    shipping_county VARCHAR(100),
    shipping_postal_code VARCHAR(20),
    shipping_country VARCHAR(100),
    shipping_phone VARCHAR(20),
    ordered_at TIMESTAMP DEFAULT NOW(),
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders.order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders.orders(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES catalog.product_variants(id),
    product_name VARCHAR(255) NOT NULL,
    variant_sku VARCHAR(100) NOT NULL,
    size VARCHAR(20),
    color VARCHAR(50),
    unit_price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders.cart_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES auth.users(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES catalog.products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, product_id)
);

CREATE TABLE orders.guest_cart_items (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    product_id INTEGER REFERENCES catalog.products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(session_id, product_id)
);

-- =====================================================
-- PERFORMANCE INDEXES
-- =====================================================

-- Catalog indexes
CREATE INDEX idx_catalog_products_category ON catalog.products(category_id);
CREATE INDEX idx_catalog_products_active ON catalog.products(is_active);
CREATE INDEX idx_catalog_products_featured ON catalog.products(is_featured);

CREATE INDEX idx_catalog_variants_product ON catalog.product_variants(product_id);
CREATE INDEX idx_catalog_variants_stock ON catalog.product_variants(stock_quantity);

CREATE INDEX idx_catalog_images_product ON catalog.product_images(product_id);
CREATE INDEX idx_catalog_images_variant ON catalog.product_images(variant_id);
CREATE INDEX idx_catalog_images_primary ON catalog.product_images(is_primary);

-- Full-text search for products
CREATE INDEX idx_catalog_products_fulltext ON catalog.products
  USING gin (to_tsvector('english', coalesce(name,'') || ' ' || coalesce(description,'')));

-- Auth indexes
CREATE INDEX idx_auth_addresses_user ON auth.user_addresses(user_id);
CREATE INDEX idx_auth_users_lower_email ON auth.users (lower(email));
CREATE INDEX idx_auth_points_transactions_user ON auth.points_transactions(user_id);
CREATE INDEX idx_auth_points_transactions_order ON auth.points_transactions(order_id);

-- Orders indexes
CREATE INDEX idx_orders_user ON orders.orders(user_id);
CREATE INDEX idx_orders_status ON orders.orders(status);
CREATE INDEX idx_orders_ordered_at ON orders.orders(ordered_at);

CREATE INDEX idx_orders_order_items_order ON orders.order_items(order_id);
CREATE INDEX idx_orders_order_items_variant ON orders.order_items(variant_id);

CREATE INDEX idx_orders_cart_variant ON orders.cart_items(variant_id);
CREATE INDEX idx_orders_guest_cart_session ON orders.guest_cart_items(session_id);

-- Partial index for active products
CREATE INDEX idx_catalog_products_active_slug ON catalog.products (slug) WHERE is_active;