
CREATE TABLE users (
   id SERIAL PRIMARY KEY,               -- Auto-incrementing user ID
   user_name VARCHAR(255) NOT NULL UNIQUE,
   email VARCHAR(255) NOT NULL UNIQUE,   -- User's email (must be unique)
   password_hash VARCHAR(255) NOT NULL, -- Hashed password
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the user was created
   updated_at TIMESTAMP DEFAULT NULL -- Timestamp when the user was last updated
);

CREATE TABLE access_tokens (
   id SERIAL PRIMARY KEY,               -- Auto-incrementing token ID
   token VARCHAR(255) NOT NULL,          -- JWT access token
   user_id INT NOT NULL,                 -- Reference to the user who owns the token
   expiry TIMESTAMP NOT NULL,           -- Expiration date of the access token
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the token was created
   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,               -- Auto-incrementing token ID
    token VARCHAR(255) NOT NULL,          -- JWT refresh token
    user_id INT NOT NULL,                 -- Reference to the user who owns the token
    expiry TIMESTAMP NOT NULL,           -- Expiration date of the refresh token
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp when the token was created
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users (user_name, email, password_hash)
VALUES ('abc', '01901901901@mailinator.com', '321dsaf');




