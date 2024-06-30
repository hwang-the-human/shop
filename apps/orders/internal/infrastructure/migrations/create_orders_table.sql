CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity INT NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'success', 'failed', 'paid')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);