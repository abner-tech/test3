

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    authors VARCHAR(255)[], -- Array to handle multiple authors
    isbn VARCHAR(13) UNIQUE NOT NULL,
    publication_date DATE,
    genre VARCHAR(50)[],
    description TEXT,
    average_rating DECIMAL(3,2) CHECK (average_rating BETWEEN 0 AND 5) DEFAULT 0.0,
    version INT NOT NULL DEFAULT 1
);