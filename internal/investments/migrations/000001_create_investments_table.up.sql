CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE investments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    account_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    symbol VARCHAR(10) NOT NULL,
    quantity DECIMAL(15,6) NOT NULL,
    purchase_price DECIMAL(15,2) NOT NULL,
    current_price DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    purchase_date TIMESTAMP NOT NULL,
    last_updated TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_investments_user_id ON investments(user_id);
CREATE INDEX idx_investments_account_id ON investments(account_id);
CREATE INDEX idx_investments_symbol ON investments(symbol);
CREATE INDEX idx_investments_status ON investments(status);

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_investments_updated_at
    BEFORE UPDATE ON investments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 