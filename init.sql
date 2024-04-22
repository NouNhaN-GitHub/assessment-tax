-- Creation of table
CREATE TYPE allowance_type AS ENUM ('donation', 'k-receipt', 'personal');

CREATE TABLE IF NOT EXISTS allowances (
    allowance_type allowance_type NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO allowances (allowance_type, amount) VALUES
('donation',0.00),
('k-receipt',50000.00),
('personal',60000.00);