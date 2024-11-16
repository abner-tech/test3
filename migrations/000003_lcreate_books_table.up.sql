-- CREATE TABLE IF NOT EXISTS reading_lists (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(100) NOT NULL,
--     description TEXT,
--     created_by INT REFERENCES users(id) ON DELETE SET NULL,
--     version INT NOT NULL DEFAULT 1
-- );

-- -- Junction table for Reading Lists and Books (many-to-many relationship)
-- CREATE TABLE IF NOT EXISTS reading_list_books (
--     reading_list_id INT REFERENCES reading_lists(id) ON DELETE CASCADE,
--     book_id INT REFERENCES books(id) ON DELETE CASCADE,
--     PRIMARY KEY (reading_list_id, book_id)
-- );

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    authors VARCHAR(255)[], -- Array to handle multiple authors
    isbn VARCHAR(13) UNIQUE NOT NULL,
    publication_date DATE,
    genre VARCHAR(50),
    description TEXT,
    average_rating DECIMAL(3,2) CHECK (average_rating BETWEEN 0 AND 5) DEFAULT 0.0
);