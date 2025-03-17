package finance

import "time"

// User - Enhanced with comprehensive profile information and related entities
type User struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	Phone          string         `json:"phone"`
	Address        Address        `json:"address"`
	DateOfBirth    time.Time      `json:"date_of_birth"`
	CreditScore    int            `json:"credit_score"`
	Role           UserRole       `json:"role"`
	Accounts       []Account      `json:"accounts"`
	Investments    []Investment   `json:"investments"`
	Loans          []Loan         `json:"loans"`
	Budgets        []Budget       `json:"budgets"`
	Subscriptions  []Subscription `json:"subscriptions"`
	KYCVerified    bool           `json:"kyc_verified"`
	RegisteredDate time.Time      `json:"registered_date"`
	LastLogin      time.Time      `json:"last_login"`
}

// Address - Common structure for locations
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	IsVerified bool   `json:"is_verified"`
}

// UserRole - Enhanced with role hierarchy
type UserRole struct {
	RoleName    string    `json:"role_name"`
	Permissions []string  `json:"permissions"`
	Level       int       `json:"level"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// AuthToken - Enhanced with security features
type AuthToken struct {
	UserID        string    `json:"user_id"`
	Token         string    `json:"token"`
	Expiry        time.Time `json:"expiry"`
	IPAddress     string    `json:"ip_address"`
	DeviceInfo    string    `json:"device_info"`
	LastUsed      time.Time `json:"last_used"`
	RefreshToken  string    `json:"refresh_token"`
	RefreshExpiry time.Time `json:"refresh_expiry"`
}

// Bank - Enhanced with regulatory information
type Bank struct {
	BankID        string    `json:"bank_id"`
	Name          string    `json:"name"`
	Country       string    `json:"country"`
	LicenseNumber string    `json:"license_number"`
	Branches      []Branch  `json:"branches"`
	Accounts      []Account `json:"accounts"`
	SwiftCode     string    `json:"swift_code"`
	RoutingNumber string    `json:"routing_number"`
	Founded       time.Time `json:"founded"`
	Website       string    `json:"website"`
}

// Branch - New structure for bank branches
type Branch struct {
	BranchID    string  `json:"branch_id"`
	BankID      string  `json:"bank_id"`
	Name        string  `json:"name"`
	Address     Address `json:"address"`
	Manager     string  `json:"manager"`
	PhoneNumber string  `json:"phone_number"`
}

// AccountType - Enum for account types
type AccountType string

const (
	Checking   AccountType = "CHECKING"
	Savings    AccountType = "SAVINGS"
	Investment AccountType = "INVESTMENT"
	Credit     AccountType = "CREDIT"
	Loan       AccountType = "LOAN"
)

// AccountStatus - Enum for account status
type AccountStatus string

const (
	Active  AccountStatus = "ACTIVE"
	Frozen  AccountStatus = "FROZEN"
	Closed  AccountStatus = "CLOSED"
	Pending AccountStatus = "PENDING"
)

// Account - Enhanced with type, status and history
type Account struct {
	AccountID        string        `json:"account_id"`
	AccountNumber    string        `json:"account_number"`
	UserID           string        `json:"user_id"`
	BankID           string        `json:"bank_id"`
	BranchID         string        `json:"branch_id"`
	Type             AccountType   `json:"type"`
	Status           AccountStatus `json:"status"`
	Balance          float64       `json:"balance"`
	AvailableBalance float64       `json:"available_balance"`
	Currency         Currency      `json:"currency"`
	InterestRate     float64       `json:"interest_rate"`
	MinimumBalance   float64       `json:"minimum_balance"`
	Transactions     []Transaction `json:"transactions"`
	OpenedDate       time.Time     `json:"opened_date"`
	LastActivity     time.Time     `json:"last_activity"`
	OverdraftLimit   float64       `json:"overdraft_limit"`
}

// TransactionType - Types of financial transactions
type TransactionType string

const (
	Deposit    TransactionType = "DEPOSIT"
	Withdrawal TransactionType = "WITHDRAWAL"
	Transfer   TransactionType = "TRANSFER"
	Payment    TransactionType = "PAYMENT"
	Fee        TransactionType = "FEE"
	Interest   TransactionType = "INTEREST"
	Adjustment TransactionType = "ADJUSTMENT"
)

// TransactionStatus - Status of transactions
type TransactionStatus string

const (
	Pending   TransactionStatus = "PENDING"
	Completed TransactionStatus = "COMPLETED"
	Failed    TransactionStatus = "FAILED"
	Reversed  TransactionStatus = "REVERSED"
	Disputed  TransactionStatus = "DISPUTED"
)

// Transaction - Enhanced with detailed information
type Transaction struct {
	TransactionID    string            `json:"transaction_id"`
	AccountID        string            `json:"account_id"`
	Type             TransactionType   `json:"type"`
	Status           TransactionStatus `json:"status"`
	Amount           float64           `json:"amount"`
	Currency         Currency          `json:"currency"`
	Date             time.Time         `json:"date"`
	Description      string            `json:"description"`
	Category         string            `json:"category"`
	MerchantName     string            `json:"merchant_name"`
	MerchantID       string            `json:"merchant_id"`
	LocationID       string            `json:"location_id"`
	SourceAccountID  string            `json:"source_account_id"`
	DestAccountID    string            `json:"destination_account_id"`
	ReferenceNumber  string            `json:"reference_number"`
	BudgetCategoryID string            `json:"budget_category_id"`
}

// Investment - Enhanced with performance tracking
type Investment struct {
	InvestmentID string          `json:"investment_id"`
	UserID       string          `json:"user_id"`
	AccountID    string          `json:"account_id"`
	Stocks       []StockHolding  `json:"stocks"`
	Bonds        []Bond          `json:"bonds"`
	Cryptos      []CryptoHolding `json:"cryptos"`
	TotalValue   float64         `json:"total_value"`
	InitialValue float64         `json:"initial_value"`
	ReturnRate   float64         `json:"return_rate"`
	LastUpdated  time.Time       `json:"last_updated"`
	RiskProfile  string          `json:"risk_profile"`
	StrategyType string          `json:"strategy_type"`
}

// StockHolding - Enhanced with purchase information
type StockHolding struct {
	HoldingID       string    `json:"holding_id"`
	StockID         string    `json:"stock_id"`
	Quantity        float64   `json:"quantity"`
	AverageBuyPrice float64   `json:"average_buy_price"`
	CurrentValue    float64   `json:"current_value"`
	PurchaseDate    time.Time `json:"purchase_date"`
	Stock           Stock     `json:"stock"`
}

// Stock - Enhanced with market information
type Stock struct {
	StockID       string    `json:"stock_id"`
	Symbol        string    `json:"symbol"`
	Name          string    `json:"name"`
	Exchange      string    `json:"exchange"`
	Price         float64   `json:"price"`
	Change        float64   `json:"change"`
	ChangePercent float64   `json:"change_percent"`
	MarketCap     float64   `json:"market_cap"`
	Volume        int64     `json:"volume"`
	PERatio       float64   `json:"pe_ratio"`
	DividendYield float64   `json:"dividend_yield"`
	Sector        string    `json:"sector"`
	Industry      string    `json:"industry"`
	LastUpdated   time.Time `json:"last_updated"`
}

// Bond - New structure for fixed income investments
type Bond struct {
	BondID           string    `json:"bond_id"`
	Name             string    `json:"name"`
	Issuer           string    `json:"issuer"`
	FaceValue        float64   `json:"face_value"`
	CouponRate       float64   `json:"coupon_rate"`
	MaturityDate     time.Time `json:"maturity_date"`
	PurchaseDate     time.Time `json:"purchase_date"`
	PurchasePrice    float64   `json:"purchase_price"`
	CurrentValue     float64   `json:"current_value"`
	PaymentFrequency string    `json:"payment_frequency"`
	Rating           string    `json:"rating"`
}

// CryptoHolding - New structure for cryptocurrency holdings
type CryptoHolding struct {
	HoldingID       string    `json:"holding_id"`
	CryptoID        string    `json:"crypto_id"`
	Quantity        float64   `json:"quantity"`
	AverageBuyPrice float64   `json:"average_buy_price"`
	CurrentValue    float64   `json:"current_value"`
	WalletAddress   string    `json:"wallet_address"`
	PurchaseDate    time.Time `json:"purchase_date"`
	Crypto          Crypto    `json:"crypto"`
}

// LoanType - Types of loans
type LoanType string

const (
	Mortgage     LoanType = "MORTGAGE"
	AutoLoan     LoanType = "AUTO_LOAN"
	StudentLoan  LoanType = "STUDENT_LOAN"
	PersonalLoan LoanType = "PERSONAL_LOAN"
	BusinessLoan LoanType = "BUSINESS_LOAN"
)

// LoanStatus - Status of loans
type LoanStatus string

const (
	Approved  LoanStatus = "APPROVED"
	Pending   LoanStatus = "PENDING"
	Rejected  LoanStatus = "REJECTED"
	Active    LoanStatus = "ACTIVE"
	Defaulted LoanStatus = "DEFAULTED"
	PaidOff   LoanStatus = "PAID_OFF"
)

// Loan - Enhanced with comprehensive loan details
type Loan struct {
	LoanID          string     `json:"loan_id"`
	UserID          string     `json:"user_id"`
	BankID          string     `json:"bank_id"`
	Type            LoanType   `json:"type"`
	Status          LoanStatus `json:"status"`
	Amount          float64    `json:"amount"`
	Interest        float64    `json:"interest"`
	Term            int        `json:"term_months"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Payments        []Payment  `json:"payments"`
	RemainingAmount float64    `json:"remaining_amount"`
	NextPaymentDate time.Time  `json:"next_payment_date"`
	Collateral      string     `json:"collateral"`
	CollateralValue float64    `json:"collateral_value"`
	APR             float64    `json:"apr"`
	LateFeesAccrued float64    `json:"late_fees_accrued"`
}

// PaymentStatus - Status of payments
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentCompleted PaymentStatus = "COMPLETED"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentLate      PaymentStatus = "LATE"
)

// Payment - Enhanced with payment status and method
type Payment struct {
	PaymentID     string        `json:"payment_id"`
	LoanID        string        `json:"loan_id"`
	AmountPaid    float64       `json:"amount_paid"`
	DatePaid      time.Time     `json:"date_paid"`
	DueDate       time.Time     `json:"due_date"`
	Status        PaymentStatus `json:"status"`
	PaymentMethod string        `json:"payment_method"`
	TransactionID string        `json:"transaction_id"`
	PrincipalPaid float64       `json:"principal_paid"`
	InterestPaid  float64       `json:"interest_paid"`
	LateFee       float64       `json:"late_fee"`
	PaymentNumber int           `json:"payment_number"`
}

// Currency - Enhanced with additional details
type Currency struct {
	Code          string         `json:"code"`
	Symbol        string         `json:"symbol"`
	Name          string         `json:"name"`
	DecimalPlaces int            `json:"decimal_places"`
	IsDigital     bool           `json:"is_digital"`
	ExchangeRates []ExchangeRate `json:"exchange_rates"`
}

// ExchangeRate - Enhanced with temporal information
type ExchangeRate struct {
	FromCurrency    string           `json:"from_currency"`
	ToCurrency      string           `json:"to_currency"`
	Rate            float64          `json:"rate"`
	Date            time.Time        `json:"date"`
	Provider        string           `json:"provider"`
	Spread          float64          `json:"spread"`
	HistoricalRates []HistoricalRate `json:"historical_rates"`
}

// HistoricalRate - New structure for tracking rate changes
type HistoricalRate struct {
	Date time.Time `json:"date"`
	Rate float64   `json:"rate"`
}

// Report - Enhanced with report types and status
type Report struct {
	ReportID        string    `json:"report_id"`
	UserID          string    `json:"user_id"`
	Title           string    `json:"title"`
	Details         string    `json:"details"`
	GeneratedDate   time.Time `json:"generated_date"`
	Type            string    `json:"type"`
	Format          string    `json:"format"`
	Status          string    `json:"status"`
	DownloadURL     string    `json:"download_url"`
	ValidUntil      time.Time `json:"valid_until"`
	RelatedEntities []string  `json:"related_entities"`
}

// Additional structures would follow with similar enhancements for:
// TaxReport, MarketData, Crypto, Expense, Budget, Savings,
// WealthManagement, Audit, Notification, Subscription
