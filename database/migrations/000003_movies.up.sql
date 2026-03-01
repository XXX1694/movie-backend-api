CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    year INT,
    rating NUMERIC(3,1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

INSERT INTO movies (title, description, year, rating) VALUES
('The Shawshank Redemption', 'Two imprisoned men bond over a number of years.', 1994, 9.3),
('The Godfather', 'The aging patriarch of an organized crime dynasty transfers control to his reluctant son.', 1972, 9.2),
('The Dark Knight', 'Batman raises the stakes in his war on crime.', 2008, 9.0);
