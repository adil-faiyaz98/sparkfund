package postgres

import (
	"fmt"

	"github.com/adil-faiyaz98/structgen/internal/loans"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) loans.LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) Create(loan *loans.Loan) error {
	return r.db.Create(loan).Error
}

func (r *loanRepository) GetByID(id uuid.UUID) (*loans.Loan, error) {
	var loan loans.Loan
	if err := r.db.First(&loan, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("loan not found: %w", err)
	}
	return &loan, nil
}

func (r *loanRepository) GetByUserID(userID uuid.UUID) ([]*loans.Loan, error) {
	var loanList []*loans.Loan
	if err := r.db.Where("user_id = ?", userID).Find(&loanList).Error; err != nil {
		return nil, fmt.Errorf("failed to get user loans: %w", err)
	}
	return loanList, nil
}

func (r *loanRepository) GetByAccountID(accountID uuid.UUID) ([]*loans.Loan, error) {
	var loanList []*loans.Loan
	if err := r.db.Where("account_id = ?", accountID).Find(&loanList).Error; err != nil {
		return nil, fmt.Errorf("failed to get account loans: %w", err)
	}
	return loanList, nil
}

func (r *loanRepository) Update(loan *loans.Loan) error {
	return r.db.Save(loan).Error
}

func (r *loanRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&loans.Loan{}, "id = ?", id).Error
}

func (r *loanRepository) CreatePayment(payment *loans.LoanPayment) error {
	return r.db.Create(payment).Error
}

func (r *loanRepository) GetPayments(loanID uuid.UUID) ([]*loans.LoanPayment, error) {
	var paymentList []*loans.LoanPayment
	if err := r.db.Where("loan_id = ?", loanID).Find(&paymentList).Error; err != nil {
		return nil, fmt.Errorf("failed to get loan payments: %w", err)
	}
	return paymentList, nil
}
