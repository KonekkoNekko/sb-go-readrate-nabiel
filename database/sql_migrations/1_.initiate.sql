-- 1_.initiate.sql
-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL, -- Increased size for hashed passwords
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by INT, -- References users.id, 0 for system or initially null
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by INT -- References users.id
);

-- Categories Table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by INT, -- References users.id
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by INT -- References users.id
);

-- Books Table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category_id INT NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    release_year INT NOT NULL,
    price INT NOT NULL,
    total_page INT NOT NULL,
    thickness VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by INT, -- References users.id
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by INT, -- References users.id
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT
);

-- Reviews Table (NEW TABLE)
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    user_id INT NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5), -- Rating from 1 to 5
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by INT, -- References users.id
    modified_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    modified_by INT, -- References users.id
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (book_id, user_id) -- A user can only review a book once
);

-- Add sample user (for testing login)
-- IMPORTANT: Replace 'securepassword' with a bcrypt hash generated from a tool or your register endpoint
-- INSERT INTO users (username, password, created_by, modified_by) VALUES ('testuser', '$2a$10$YourGeneratedBcryptHashHere...', 0, 0);
-- Example hash for 'password': $2a$10$e.g.Y6W4lH7uX9P4n2g5a.Vf8o7j9k0l1m2n3o4p5q6r7s8t9u0v1w2x3y4z5A6B7C
-- You should register a user via the /register endpoint to get a valid hash.