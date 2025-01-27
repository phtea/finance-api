-- +goose Up
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,             -- Transaction ID
    sender_id INT NOT NULL,            -- Sender ID
    receiver_id INT NOT NULL,          -- Receiver ID
    amount DECIMAL(10, 2) NOT NULL,    -- Transaction amount
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Timestamp of the transaction
    FOREIGN KEY (sender_id) REFERENCES users(id),  -- Foreign key linking sender to users table
    FOREIGN KEY (receiver_id) REFERENCES users(id) -- Foreign key linking receiver to users table
);

-- +goose Down
DROP TABLE transactions;
