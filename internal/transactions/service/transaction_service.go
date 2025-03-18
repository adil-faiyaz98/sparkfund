package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/transactions"
	"github.com/google/uuid"
)

type transactionService struct {
	repo transactions.TransactionRepository
}

func NewTransactionService(repo transactions.TransactionRepository) transactions.TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(ctx context.Context, transaction *transactions.Transaction) error {
	// Validate transaction type
	if !isValidTransactionType(transaction.Type) {
		return fmt.Errorf("invalid transaction type: %s", transaction.Type)
	}

	// Validate amount
	if transaction.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// Validate currency
	if !isValidCurrency(transaction.Currency) {
		return fmt.Errorf("invalid currency: %s", transaction.Currency)
	}

	// Set initial status
	transaction.Status = transactions.TransactionStatusPending

	return s.repo.Create(transaction)
}

func (s *transactionService) GetTransaction(ctx context.Context, id uuid.UUID) (*transactions.Transaction, error) {
	return s.repo.GetByID(id)
}

func (s *transactionService) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]*transactions.Transaction, error) {
	return s.repo.GetByUserID(userID)
}

func (s *transactionService) GetAccountTransactions(ctx context.Context, accountID uuid.UUID) ([]*transactions.Transaction, error) {
	return s.repo.GetByAccountID(accountID)
}

func (s *transactionService) UpdateTransactionStatus(ctx context.Context, id uuid.UUID, status transactions.TransactionStatus, err error) error {
	transaction, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	transaction.Status = status

	now := time.Now()
	switch status {
	case transactions.TransactionStatusCompleted:
		transaction.CompletedAt = &now
	case transactions.TransactionStatusFailed:
		transaction.FailedAt = &now
		if err != nil {
			transaction.Metadata = fmt.Sprintf(`{"error": "%s"}`, err.Error())
		}
	case transactions.TransactionStatusCancelled:
		transaction.CancelledAt = &now
	}

	return s.repo.Update(transaction)
}

func (s *transactionService) DeleteTransaction(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func isValidTransactionType(transactionType transactions.TransactionType) bool {
	switch transactionType {
	case transactions.TransactionTypeDeposit,
		transactions.TransactionTypeWithdrawal,
		transactions.TransactionTypeTransfer,
		transactions.TransactionTypePayment,
		transactions.TransactionTypeInterest,
		transactions.TransactionTypeFee:
		return true
	default:
		return false
	}
}

func isValidCurrency(currency string) bool {
	// In a real application, this would check against a list of valid ISO 4217 currency codes
	return len(currency) == 3
}
