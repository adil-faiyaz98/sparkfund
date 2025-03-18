package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/loans"
	"github.com/google/uuid"
)

type loanService struct {
	repo loans.LoanRepository
}

func NewLoanService(repo loans.LoanRepository) loans.LoanService {
	return &loanService{repo: repo}
}

func (s *loanService) CreateLoan(ctx context.Context, loan *loans.Loan) error {
	// Validate loan type
	if !isValidLoanType(loan.Type) {
		return fmt.Errorf("invalid loan type: %s", loan.Type)
	}

	// Calculate loan details
	monthlyPayment, totalInterest, totalAmount := calculateLoanDetails(
		loan.Amount,
		loan.InterestRate,
		loan.TermMonths,
	)

	loan.MonthlyPayment = monthlyPayment
	loan.TotalInterest = totalInterest
	loan.TotalAmount = totalAmount
	loan.Status = loans.LoanStatusPending

	return s.repo.Create(loan)
}

func (s *loanService) GetLoan(ctx context.Context, id uuid.UUID) (*loans.Loan, error) {
	return s.repo.GetByID(id)
}

func (s *loanService) GetUserLoans(ctx context.Context, userID uuid.UUID) ([]*loans.Loan, error) {
	return s.repo.GetByUserID(userID)
}

func (s *loanService) GetAccountLoans(ctx context.Context, accountID uuid.UUID) ([]*loans.Loan, error) {
	return s.repo.GetByAccountID(accountID)
}

func (s *loanService) UpdateLoanStatus(ctx context.Context, id uuid.UUID, status loans.LoanStatus, notes string) error {
	loan, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	loan.Status = status
	loan.Notes = notes

	now := time.Now()
	switch status {
	case loans.LoanStatusApproved:
		loan.ApprovedAt = &now
	case loans.LoanStatusRejected:
		loan.RejectedAt = &now
	case loans.LoanStatusPaid:
		loan.PaidAt = &now
	case loans.LoanStatusDefaulted:
		loan.DefaultedAt = &now
	}

	return s.repo.Update(loan)
}

func (s *loanService) MakePayment(ctx context.Context, loanID uuid.UUID, amount float64) error {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		return fmt.Errorf("loan not found: %w", err)
	}

	if loan.Status != loans.LoanStatusActive {
		return fmt.Errorf("loan is not active")
	}

	payment := &loans.LoanPayment{
		ID:          uuid.New(),
		LoanID:      loanID,
		Amount:      amount,
		PaymentDate: time.Now(),
	}

	// Calculate principal and interest portions
	principal, interest := calculatePaymentBreakdown(amount, loan.MonthlyPayment)
	payment.Principal = principal
	payment.Interest = interest

	return s.repo.CreatePayment(payment)
}

func (s *loanService) GetLoanPayments(ctx context.Context, loanID uuid.UUID) ([]*loans.LoanPayment, error) {
	return s.repo.GetPayments(loanID)
}

func (s *loanService) DeleteLoan(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func isValidLoanType(loanType loans.LoanType) bool {
	switch loanType {
	case loans.LoanTypePersonal,
		loans.LoanTypeMortgage,
		loans.LoanTypeBusiness,
		loans.LoanTypeEducation:
		return true
	default:
		return false
	}
}

func calculateLoanDetails(amount, interestRate float64, termMonths int) (monthlyPayment, totalInterest, totalAmount float64) {
	// Monthly interest rate
	r := interestRate / 100 / 12
	// Number of payments
	n := float64(termMonths)

	// Monthly payment formula: P = L[r(1 + r)^n]/[(1 + r)^n - 1]
	monthlyPayment = amount * (r * pow(1+r, n)) / (pow(1+r, n) - 1)
	totalAmount = monthlyPayment * n
	totalInterest = totalAmount - amount

	return monthlyPayment, totalInterest, totalAmount
}

func calculatePaymentBreakdown(payment, monthlyPayment float64) (principal, interest float64) {
	if payment >= monthlyPayment {
		principal = monthlyPayment
		interest = payment - monthlyPayment
	} else {
		principal = payment
		interest = 0
	}
	return principal, interest
}

func pow(x, y float64) float64 {
	result := 1.0
	for i := 0; i < int(y); i++ {
		result *= x
	}
	return result
}
