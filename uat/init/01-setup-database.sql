-- Test database schema and data for pg-goer UAT
-- This creates a realistic e-commerce-like schema to test all features

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    is_verified BOOLEAN DEFAULT false
);

-- Create orders table with foreign key to users
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    shipping_address TEXT,
    notes TEXT
);

-- Create order_items table with foreign key to orders
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) GENERATED ALWAYS AS (quantity * unit_price) STORED
);

-- Create categories table for additional relationships
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id INTEGER REFERENCES categories(id)
);

-- Create products table with foreign key to categories
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    in_stock BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert test data
INSERT INTO users (email, first_name, last_name, status, is_verified) VALUES
('john.doe@example.com', 'John', 'Doe', 'active', true),
('jane.smith@example.com', 'Jane', 'Smith', 'active', true),
('bob.wilson@example.com', 'Bob', 'Wilson', 'active', false),
('alice.brown@example.com', 'Alice', 'Brown', 'inactive', true),
('charlie.davis@example.com', 'Charlie', 'Davis', 'active', true),
('diana.miller@example.com', 'Diana', 'Miller', 'active', true),
('edward.jones@example.com', 'Edward', 'Jones', 'active', false),
('fiona.garcia@example.com', 'Fiona', 'Garcia', 'active', true),
('george.martinez@example.com', 'George', 'Martinez', 'inactive', false),
('helen.anderson@example.com', 'Helen', 'Anderson', 'active', true);

INSERT INTO categories (name, description, parent_id) VALUES
('Electronics', 'Electronic devices and accessories', NULL),
('Computers', 'Desktop and laptop computers', 1),
('Phones', 'Mobile phones and smartphones', 1),
('Books', 'Physical and digital books', NULL),
('Fiction', 'Fiction books and novels', 4),
('Non-Fiction', 'Educational and reference books', 4);

INSERT INTO products (name, description, price, category_id, in_stock) VALUES
('Gaming Laptop', 'High-performance gaming laptop', 1299.99, 2, true),
('Smartphone Pro', 'Latest flagship smartphone', 899.99, 3, true),
('Wireless Headphones', 'Noise-canceling wireless headphones', 249.99, 1, true),
('Programming Book', 'Learn Go programming language', 49.99, 6, true),
('Science Fiction Novel', 'Bestselling sci-fi adventure', 14.99, 5, true),
('Desktop Computer', 'Powerful desktop workstation', 1899.99, 2, false),
('Tablet Device', 'Lightweight tablet for productivity', 599.99, 1, true),
('Mystery Novel', 'Thrilling mystery story', 12.99, 5, true);

INSERT INTO orders (user_id, total, status, shipping_address, notes) VALUES
(1, 1299.99, 'completed', '123 Main St, City, State 12345', 'Express delivery requested'),
(2, 899.99, 'completed', '456 Oak Ave, Town, State 67890', NULL),
(1, 264.98, 'shipped', '123 Main St, City, State 12345', 'Gift wrapping requested'),
(3, 49.99, 'completed', '789 Pine Rd, Village, State 13579', NULL),
(4, 27.98, 'pending', '321 Elm St, Borough, State 24680', 'Hold for pickup'),
(5, 599.99, 'processing', '654 Cedar Ln, Hamlet, State 97531', NULL),
(2, 1899.99, 'completed', '456 Oak Ave, Town, State 67890', 'Business purchase'),
(6, 12.99, 'completed', '987 Birch Dr, Township, State 86420', NULL),
(7, 149.99, 'cancelled', '147 Maple Way, County, State 75319', 'Customer requested cancellation'),
(1, 49.99, 'completed', '123 Main St, City, State 12345', NULL);

INSERT INTO order_items (order_id, product_name, quantity, unit_price) VALUES
(1, 'Gaming Laptop', 1, 1299.99),
(2, 'Smartphone Pro', 1, 899.99),
(3, 'Wireless Headphones', 1, 249.99),
(3, 'Programming Book', 1, 49.99),
(3, 'Science Fiction Novel', 1, 14.99),
(4, 'Programming Book', 1, 49.99),
(5, 'Mystery Novel', 1, 12.99),
(5, 'Science Fiction Novel', 1, 14.99),
(6, 'Tablet Device', 1, 599.99),
(7, 'Desktop Computer', 1, 1899.99),
(8, 'Mystery Novel', 1, 12.99),
(9, 'Wireless Headphones', 1, 149.99),
(10, 'Programming Book', 1, 49.99);

-- Update statistics to ensure row counts are accurate
ANALYZE users;
ANALYZE orders;
ANALYZE order_items;
ANALYZE categories;
ANALYZE products;

-- Create some indexes for realism
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_date ON orders(order_date);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_products_category_id ON products(category_id);

-- Enable useful PostgreSQL extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create audit table for demonstration
CREATE TABLE audit_log (
    id SERIAL PRIMARY KEY,
    table_name VARCHAR(50) NOT NULL,
    operation VARCHAR(10) NOT NULL,
    user_name VARCHAR(100) DEFAULT current_user,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    old_values JSONB,
    new_values JSONB
);

-- Create trigger function for auditing
CREATE OR REPLACE FUNCTION audit_trigger()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        INSERT INTO audit_log (table_name, operation, old_values)
        VALUES (TG_TABLE_NAME, TG_OP, row_to_json(OLD));
        RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_log (table_name, operation, old_values, new_values)
        VALUES (TG_TABLE_NAME, TG_OP, row_to_json(OLD), row_to_json(NEW));
        RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        INSERT INTO audit_log (table_name, operation, new_values)
        VALUES (TG_TABLE_NAME, TG_OP, row_to_json(NEW));
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create audit triggers
CREATE TRIGGER users_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE FUNCTION audit_trigger();

CREATE TRIGGER orders_audit_trigger
    AFTER INSERT OR UPDATE OR DELETE ON orders
    FOR EACH ROW EXECUTE FUNCTION audit_trigger();

-- Create views for demonstration
CREATE VIEW active_users AS
SELECT 
    id,
    email,
    first_name || ' ' || last_name AS full_name,
    created_at,
    is_verified
FROM users 
WHERE status = 'active';

CREATE VIEW order_summary AS
SELECT 
    o.id,
    u.email,
    o.total,
    o.status,
    o.order_date,
    COUNT(oi.id) as item_count
FROM orders o
JOIN users u ON o.user_id = u.id
LEFT JOIN order_items oi ON o.id = oi.order_id
GROUP BY o.id, u.email, o.total, o.status, o.order_date;

-- Create some additional indexes with different types
CREATE INDEX CONCURRENTLY idx_users_status_verified ON users(status, is_verified);
CREATE INDEX idx_products_price_btree ON products USING btree(price);
CREATE INDEX idx_orders_total_range ON orders(total) WHERE total > 100;

-- Create a function for demonstration
CREATE OR REPLACE FUNCTION get_user_order_count(user_email VARCHAR)
RETURNS INTEGER AS $$
DECLARE
    order_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO order_count
    FROM orders o
    JOIN users u ON o.user_id = u.id
    WHERE u.email = user_email;
    
    RETURN order_count;
END;
$$ LANGUAGE plpgsql;

-- Create a more complex trigger for order validation
CREATE OR REPLACE FUNCTION validate_order()
RETURNS TRIGGER AS $$
BEGIN
    -- Ensure order total is positive
    IF NEW.total <= 0 THEN
        RAISE EXCEPTION 'Order total must be positive, got %', NEW.total;
    END IF;
    
    -- Ensure user exists and is active
    IF NOT EXISTS (
        SELECT 1 FROM users 
        WHERE id = NEW.user_id AND status = 'active'
    ) THEN
        RAISE EXCEPTION 'Cannot create order for inactive or non-existent user %', NEW.user_id;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER order_validation_trigger
    BEFORE INSERT OR UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION validate_order();