-- Integration test schema - simplified for fast CI execution
-- Focus on covering all feature areas with minimal data

-- Users table with various column types and constraints
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL,
    full_name TEXT,
    age INTEGER CHECK (age >= 0),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB
);

-- Posts table with foreign key to users
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    published_at TIMESTAMP,
    view_count INTEGER DEFAULT 0,
    is_published BOOLEAN DEFAULT false
);

-- Comments table with foreign key to both users and posts
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    parent_comment_id INTEGER REFERENCES comments(id)
);

-- Categories table for self-referencing relationship
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id INTEGER REFERENCES categories(id),
    sort_order INTEGER DEFAULT 0
);

-- Junction table for many-to-many relationship
CREATE TABLE post_categories (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);

-- Insert minimal test data
INSERT INTO users (username, email, full_name, age, is_active) VALUES
('admin', 'admin@example.com', 'System Administrator', 35, true),
('john_doe', 'john@example.com', 'John Doe', 28, true),
('jane_smith', 'jane@example.com', 'Jane Smith', 32, false);

INSERT INTO categories (name, slug, description, parent_id, sort_order) VALUES
('Technology', 'technology', 'Technology related posts', NULL, 1),
('Programming', 'programming', 'Programming tutorials and tips', 1, 1),
('Web Development', 'web-dev', 'Web development content', 2, 1),
('Lifestyle', 'lifestyle', 'Lifestyle content', NULL, 2);

INSERT INTO posts (title, content, author_id, published_at, view_count, is_published) VALUES
('Getting Started with Go', 'Learn the basics of Go programming...', 1, CURRENT_TIMESTAMP, 150, true),
('Database Design Tips', 'Best practices for database design...', 2, CURRENT_TIMESTAMP, 89, true),
('Draft Post', 'This is a draft post...', 1, NULL, 0, false);

INSERT INTO comments (post_id, author_id, content, parent_comment_id) VALUES
(1, 2, 'Great tutorial! Very helpful.', NULL),
(1, 3, 'Thanks for sharing this.', NULL),
(1, 1, 'Glad you found it useful!', 1),
(2, 1, 'Excellent points on normalization.', NULL);

INSERT INTO post_categories (post_id, category_id) VALUES
(1, 2), -- Programming
(1, 3), -- Web Development
(2, 1), -- Technology
(2, 2); -- Programming

-- Update statistics for accurate row counts
ANALYZE;

-- Create some indexes for realism
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_published_at ON posts(published_at);
CREATE INDEX idx_comments_post_id ON comments(post_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);