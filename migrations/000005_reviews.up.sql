--script to create the reviews table
CREATE TABLE IF NOT EXISTS reviews (
    id bigserial PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE CASCADE, --foreign key
    user_name VARCHAR(255), --name of user who posted the rating message and value
    rating INT CHECK(rating BETWEEN 1 AND 5), --rating value between 1 and 5
    review_text TEXT, --rating message
    helpful_count INT DEFAULT 0, --count to see amount of persons who found this rating helpful
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    version integer NOT NULL DEFAULT 1
);