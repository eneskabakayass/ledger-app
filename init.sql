-- Create users table
CREATE TABLE IF NOT EXISTS users (
                                     id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Insert sample users
INSERT INTO users (name, password_hash, is_admin) VALUES ('John Doe', 'hashed_password_123', FALSE);
INSERT INTO users (name, password_hash, is_admin) VALUES ('Jane Smith', 'hashed_password_456', TRUE);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
                                            id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
                                            user_id BIGINT UNSIGNED NOT NULL,
                                            amount DOUBLE NOT NULL,
                                            transaction_time TIMESTAMP NOT NULL,
                                            sender_id BIGINT UNSIGNED,
                                            receiver_id BIGINT UNSIGNED,
                                            CONSTRAINT fk_users_transactions FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_transactions_sender FOREIGN KEY (sender_id) REFERENCES users(id),
    CONSTRAINT fk_transactions_receiver FOREIGN KEY (receiver_id) REFERENCES users(id),
    INDEX idx_transactions_sender_id (sender_id),
    INDEX idx_transactions_receiver_id (receiver_id)
    );

-- Example of inserting a transaction
INSERT INTO transactions (user_id, amount, transaction_time, sender_id, receiver_id)
VALUES (1, 100.50, CURRENT_TIMESTAMP, NULL, 2);