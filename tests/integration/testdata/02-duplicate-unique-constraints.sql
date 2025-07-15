-- Test case for duplicate unique constraints on the same column
-- This reproduces the issue reported in GitHub Issue #6

-- Create a table with duplicate unique constraints
CREATE TABLE magic_code_rate_limits (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email varchar(255) NOT NULL,
    attempts integer DEFAULT 0,
    last_attempt timestamp,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);

-- Create the first unique constraint
ALTER TABLE magic_code_rate_limits ADD CONSTRAINT magic_code_rate_limits_email_key UNIQUE (email);

-- Create a duplicate unique constraint on the same column
-- This is what causes the issue in the original GetColumns query
ALTER TABLE magic_code_rate_limits ADD CONSTRAINT magic_code_rate_limits_email_unique UNIQUE (email);

-- Create another table with multiple unique constraints on different columns (should work fine)
CREATE TABLE test_multiple_uniques (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    CONSTRAINT test_code_unique UNIQUE (code),
    CONSTRAINT test_name_unique UNIQUE (name),
    CONSTRAINT test_email_unique UNIQUE (email)
);

-- Insert test data
INSERT INTO magic_code_rate_limits (email, attempts) VALUES
('test@example.com', 1),
('user@example.com', 3);

INSERT INTO test_multiple_uniques (code, name, email) VALUES
('ABC123', 'Test Name', 'test@example.com'),
('XYZ789', 'Another Name', 'another@example.com');

-- Update statistics
ANALYZE;