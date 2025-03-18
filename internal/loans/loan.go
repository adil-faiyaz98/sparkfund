package loans

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LoanType string

const (
	LoanTypePersonal  LoanType = "PERSONAL"
	LoanTypeMortgage  LoanType = "MORTGAGE"
	LoanTypeBusiness  LoanType = "BUSINESS"
	LoanTypeEducation LoanType = "EDUCATION"
)

type LoanStatus string

const (
	LoanStatusPending   LoanStatus = "PENDING"
	LoanStatusApproved  LoanStatus = "APPROVED"
	LoanStatusRejected  LoanStatus = "REJECTED"
	LoanStatusActive    LoanStatus = "ACTIVE"
	LoanStatusPaid      LoanStatus = "PAID"
	LoanStatusDefaulted LoanStatus = "DEFAULTED"
)

type Loan struct {
	ID             uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserID         uuid.UUID  `json:"user_id" gorm:"type:uuid;not null"`
	AccountID      uuid.UUID  `json:"account_id" gorm:"type:uuid;not null"`
	Type           LoanType   `json:"type" gorm:"type:varchar(20);not null"`
	Amount         float64    `json:"amount" gorm:"type:decimal(15,2);not null"`
	TermMonths     int        `json:"term_months" gorm:"type:integer;not null"`
	Purpose        string     `json:"purpose" gorm:"type:text"`
	InterestRate   float64    `json:"interest_rate" gorm:"type:decimal(5,2);not null"`
	Status         LoanStatus `json:"status" gorm:"type:varchar(20);not null;default:'PENDING'"`
	MonthlyPayment float64    `json:"monthly_payment" gorm:"type:decimal(15,2);not null"`
	TotalInterest  float64    `json:"total_interest" gorm:"type:decimal(15,2);not null"`
	TotalAmount    float64    `json:"total_amount" gorm:"type:decimal(15,2);not null"`
	Notes          string     `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	ApprovedAt     *time.Time `json:"approved_at"`
	RejectedAt     *time.Time `json:"rejected_at"`
	PaidAt         *time.Time `json:"paid_at"`
	DefaultedAt    *time.Time `json:"defaulted_at"`
}

type LoanPayment struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	LoanID      uuid.UUID `json:"loan_id" gorm:"type:uuid;not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(15,2);not null"`
	PaymentDate time.Time `json:"payment_date" gorm:"type:timestamp;not null"`
	Principal   float64   `json:"principal" gorm:"type:decimal(15,2);not null"`
	Interest    float64   `json:"interest" gorm:"type:decimal(15,2);not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoanRepository interface {
	Create(loan *Loan) error
	GetByID(id uuid.UUID) (*Loan, error)
	GetByUserID(userID uuid.UUID) ([]*Loan, error)
	GetByAccountID(accountID uuid.UUID) ([]*Loan, error)
	Update(loan *Loan) error
	Delete(id uuid.UUID) error
	CreatePayment(payment *LoanPayment) error
	GetPayments(loanID uuid.UUID) ([]*LoanPayment, error)
}

type LoanService interface {
	CreateLoan(ctx context.Context, loan *Loan) error
	GetLoan(ctx context.Context, id uuid.UUID) (*Loan, error)
	GetUserLoans(ctx context.Context, userID uuid.UUID) ([]*Loan, error)
	GetAccountLoans(ctx context.Context, accountID uuid.UUID) ([]*Loan, error)
	UpdateLoanStatus(ctx context.Context, id uuid.UUID, status LoanStatus, notes string) error
	MakePayment(ctx context.Context, loanID uuid.UUID, amount float64) error
	GetLoanPayments(ctx context.Context, loanID uuid.UUID) ([]*LoanPayment, error)
	DeleteLoan(ctx context.Context, id uuid.UUID) error
}
