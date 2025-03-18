DROP TRIGGER IF EXISTS update_loan_payments_updated_at ON loan_payments;
DROP TRIGGER IF EXISTS update_loans_updated_at ON loans;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS loan_payments;
DROP TABLE IF EXISTS loans; 