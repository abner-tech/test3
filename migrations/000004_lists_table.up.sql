
CREATE TABLE IF NOT EXISTS reading_lists (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by INT REFERENCES users(id) ON DELETE SET NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1
);

-- Junction table for Reading Lists and Books (many-to-many relationship)
CREATE TABLE IF NOT EXISTS reading_list_books (
    reading_list_id INT REFERENCES reading_lists(id) ON DELETE CASCADE,
    book_id INT REFERENCES books(id) ON DELETE CASCADE,
    status VARCHAR(20) CHECK (status IN ('currently reading', 'completed')) NOT NULL DEFAULT 'currently reading',
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version INT NOT NULL DEFAULT 1,
    PRIMARY KEY (reading_list_id, book_id)
);