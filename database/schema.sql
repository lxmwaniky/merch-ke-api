-- Swags Store Ke - Complete Database Schema
-- Run these commands in your PostgreSQL terminal

-- =====================================================
-- 1. CATEGORIES TABLE (Dynamic, Hierarchical)
-- =====================================================
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    parent_id INTEGER REFERENCES categories(id),
    image_url VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 2. PRODUCTS TABLE (Base Products)
-- =====================================================
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    short_description VARCHAR(500),
    category_id INTEGER REFERENCES categories(id),
    base_price DECIMAL(10,2) NOT NULL,
    sku_prefix VARCHAR(20), -- e.g., 'GOPHER-TEE'
    is_active BOOLEAN DEFAULT true,
    is_featured BOOLEAN DEFAULT false,
    weight DECIMAL(8,2), -- for shipping calculations
    dimensions VARCHAR(100), -- "L x W x H"
    meta_title VARCHAR(255), -- SEO
    meta_description VARCHAR(500), -- SEO
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 3. PRODUCT VARIANTS (Size, Color, etc.)
-- =====================================================
CREATE TABLE product_variants (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
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

-- =====================================================
-- 4. PRODUCT IMAGES (Multiple Images per Product/Variant)
-- =====================================================
CREATE TABLE product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES product_variants(id) ON DELETE SET NULL,
    image_url VARCHAR(500) NOT NULL,
    image_path VARCHAR(500) NOT NULL, -- local path for your CDN
    image_type VARCHAR(50) DEFAULT 'gallery', -- 'main', 'gallery', 'thumbnail', 'zoom'
    alt_text VARCHAR(255),
    display_order INTEGER DEFAULT 1,
    file_size INTEGER, -- in bytes
    width INTEGER, -- image width in pixels
    height INTEGER, -- image height in pixels
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 5. USERS TABLE (Customers & Admins)
-- =====================================================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    role VARCHAR(20) DEFAULT 'customer', -- 'customer', 'admin', 'super_admin'
    is_active BOOLEAN DEFAULT true,
    email_verified BOOLEAN DEFAULT false,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 6. USER ADDRESSES (Shipping)
-- =====================================================
CREATE TABLE user_addresses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) DEFAULT 'shipping', -- 'shipping', 'billing'
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

-- =====================================================
-- 7. SHOPPING CART (User Cart)
-- =====================================================
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, variant_id)
);

-- =====================================================
-- 7B. GUEST CART (Session-based)
-- =====================================================
CREATE TABLE guest_cart_items (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL, -- Session identifier for guest users
    variant_id INTEGER REFERENCES product_variants(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(session_id, variant_id)
);

-- =====================================================
-- 7C. USER POINTS SYSTEM
-- =====================================================
CREATE TABLE user_points (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    points_balance INTEGER DEFAULT 0,
    total_earned INTEGER DEFAULT 0,
    total_spent INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 7D. POINTS TRANSACTIONS
-- =====================================================
CREATE TABLE points_transactions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    order_id INTEGER REFERENCES orders(id),
    transaction_type VARCHAR(20) NOT NULL, -- 'earned', 'spent', 'expired'
    points INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 8. ORDERS
-- =====================================================
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'processing', 'shipped', 'delivered', 'cancelled'
    payment_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'paid', 'failed', 'refunded'
    payment_method VARCHAR(50), -- 'mpesa', 'card', 'bank_transfer'
    payment_reference VARCHAR(100),
    
    -- Pricing
    subtotal DECIMAL(10,2) NOT NULL,
    tax_amount DECIMAL(10,2) DEFAULT 0.00,
    shipping_amount DECIMAL(10,2) DEFAULT 0.00,
    discount_amount DECIMAL(10,2) DEFAULT 0.00,
    total_amount DECIMAL(10,2) NOT NULL,
    
    -- Shipping Info (stored at time of order)
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
    
    -- Timestamps
    ordered_at TIMESTAMP DEFAULT NOW(),
    shipped_at TIMESTAMP,
    delivered_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 9. ORDER ITEMS (What was bought)  
-- =====================================================
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    variant_id INTEGER REFERENCES product_variants(id),
    
    -- Store product info at time of purchase (in case product changes)
    product_name VARCHAR(255) NOT NULL,
    variant_sku VARCHAR(100) NOT NULL,
    size VARCHAR(20),
    color VARCHAR(50),
    unit_price DECIMAL(10,2) NOT NULL,
    quantity INTEGER NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- =====================================================
-- 10. INDEXES FOR PERFORMANCE
-- =====================================================

-- Products indexes
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_active ON products(is_active);
CREATE INDEX idx_products_featured ON products(is_featured);
CREATE INDEX idx_products_slug ON products(slug);

-- Variants indexes
CREATE INDEX idx_variants_product ON product_variants(product_id);
CREATE INDEX idx_variants_sku ON product_variants(sku);
CREATE INDEX idx_variants_active ON product_variants(is_active);

-- Images indexes
CREATE INDEX idx_images_product ON product_images(product_id);
CREATE INDEX idx_images_variant ON product_images(variant_id);
CREATE INDEX idx_images_primary ON product_images(is_primary);

-- Orders indexes
CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_number ON orders(order_number);
CREATE INDEX idx_orders_date ON orders(ordered_at);

-- Cart indexes
CREATE INDEX idx_cart_user ON cart_items(user_id);

-- =====================================================
-- 11. SAMPLE DATA (Optional - for testing)
-- =====================================================

-- Insert sample categories
INSERT INTO categories (name, slug, description) VALUES 
('Clothing', 'clothing', 'T-shirts, hoodies, and apparel'),
('Accessories', 'accessories', 'Mugs, bags, and accessories'),
('Tech Gear', 'tech-gear', 'Gadgets and tech accessories'),
('Stickers', 'stickers', 'Laptop and bumper stickers');

-- Insert subcategories
INSERT INTO categories (name, slug, description, parent_id) VALUES 
('T-Shirts', 't-shirts', 'Graphic and plain t-shirts', 1),
('Hoodies', 'hoodies', 'Zip-up and pullover hoodies', 1),
('Mugs', 'mugs', 'Coffee mugs and drinkware', 2),
('Bags', 'bags', 'Backpacks and tote bags', 2);

-- Insert sample products
INSERT INTO products (name, slug, description, category_id, base_price, sku_prefix) VALUES 
('Go Gopher T-Shirt', 'go-gopher-tshirt', 'Official Go programming language mascot t-shirt', 5, 1500.00, 'GO-TEE'),
('Linux Penguin Mug', 'linux-penguin-mug', 'Ceramic mug with Tux the Linux penguin', 7, 800.00, 'LNX-MUG'),
('Python Developer Hoodie', 'python-dev-hoodie', 'Comfortable hoodie for Python developers', 6, 3500.00, 'PY-HOOD');

-- Insert sample variants
INSERT INTO product_variants (product_id, sku, size, color, stock_quantity) VALUES 
(1, 'GO-TEE-S-BLK', 'S', 'Black', 10),
(1, 'GO-TEE-M-BLK', 'M', 'Black', 15),
(1, 'GO-TEE-L-BLK', 'L', 'Black', 8),
(1, 'GO-TEE-M-WHT', 'M', 'White', 12),
(2, 'LNX-MUG-REG-WHT', 'Regular', 'White', 25),
(3, 'PY-HOOD-M-NVY', 'M', 'Navy', 5),
(3, 'PY-HOOD-L-NVY', 'L', 'Navy', 7);

-- =====================================================
-- END OF SCHEMA
-- =====================================================