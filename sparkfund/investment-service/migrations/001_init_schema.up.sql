CREATE TABLE IF NOT EXISTS portfolios (
    portfolio_id SERIAL PRIMARY KEY,
    client_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS investments (
    investment_id VARCHAR(36) PRIMARY KEY,
    client_id VARCHAR(36) NOT NULL,
    portfolio_id VARCHAR(36) NOT NULL,
    type VARCHAR(255) NOT NULL,
    amount DECIMAL NOT NULL,
    purchase_date DATE NOT NULL
);