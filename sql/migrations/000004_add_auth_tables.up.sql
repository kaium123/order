
CREATE TABLE users (
   id SERIAL PRIMARY KEY,
   user_name VARCHAR(255) NOT NULL UNIQUE,
   email VARCHAR(255) NOT NULL UNIQUE,
   password_hash VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE access_tokens (
   id SERIAL PRIMARY KEY,
   token VARCHAR(255) NOT NULL,
   user_id INT NOT NULL,
   expiry TIMESTAMP NOT NULL,
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    user_id INT NOT NULL,
    expiry TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

INSERT INTO users (user_name, email, password_hash)
VALUES ('01901901901@mailinator.com', '01901901901@mailinator.com', '$2a$10$QLWmBwCTIa4HEOEONaIK7uubsyk2vwZIqqvZizKjhiJvQGSPm8qcu');




