-- 1_.initiate.sql

-- migrate Up
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

-- ----------------------------------------------------------------------
-- START: Initial Data Inserts
-- ----------------------------------------------------------------------

-- Insert Sample User (for testing login)
-- IMPORTANT:
-- DO NOT insert plaintext passwords here.
-- Instead, use the /register endpoint to create users and get a bcrypt hash.
-- Once you have a hash, you can uncomment and insert, e.g.:
-- INSERT INTO users (username, password, created_by, modified_by) VALUES ('testuser', '$2a$10$YOUR_BCRYPT_HASH_HERE.e.g.Y6W4lH7uX9P4n2g5a', 0, 0);

-- Insert Sample Categories
INSERT INTO categories (name, created_by, modified_by) VALUES
('Fiction', 0, 0),
('Science Fiction', 0, 0),
('Fantasy', 0, 0),
('Non-Fiction', 0, 0);

-- Insert Sample Books
-- IMPORTANT: Ensure category_id values match existing IDs in your categories table.
-- If categories are inserted above, their IDs will likely be 1, 2, 3, 4 respectively.
INSERT INTO books (title, category_id, description, image_url, release_year, price, total_page, thickness, created_by, modified_by) VALUES
('The Hitchhiker''s Guide to the Galaxy', 2, 'A comedic science fiction series.', 'http://example.com/hitchhiker.jpg', 1979, 150000, 193, 'tipis', 0, 0),
('Lord of the Rings', 3, 'An epic fantasy adventure.', 'http://example.com/lotr.jpg', 1954, 250000, 1178, 'tebal', 0, 0),
('Sapiens: A Brief History of Humankind', 4, 'A non-fiction book exploring the history of humanity.', 'http://example.com/sapiens.jpg', 2011, 200000, 443, 'tebal', 0, 0),
('1984', 1, 'A dystopian social science fiction novel.', 'http://example.com/1984.jpg', 1949, 120000, 328, 'tebal', 0, 0);

-- ----------------------------------------------------------------------
-- END: Initial Data Inserts
-- ----------------------------------------------------------------------


-- migrate Down
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;