CREATE TABLE wikis (
                       id SERIAL PRIMARY KEY,
                       topic VARCHAR(100) NOT NULL,
                       description TEXT,
                       created_at TIMESTAMP NOT NULL,
                       updated_at TIMESTAMP NOT NULL
);