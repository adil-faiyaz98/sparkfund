package handler

import (
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/service"
	// We would need to generate the protobuf code and import it here
	// pb "github.com/adil-faiyaz98/money-pulse/api/accounts"
)

// AccountHandler handles gRPC requests for the account service
type AccountHandler struct {
	service *service.AccountService
	// When protobuf is generated, uncomment this line:
	// pb.UnimplementedAccountServiceServer
}

// NewAccountHandler creates a new gRPC handler for accounts
func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

/*
// Note: These method implementations depend on the generated protobuf code.
// They are commented out until we set up the protobuf generation.

// CreateAccount implements the gRPC CreateAccount method
func (h *AccountHandler) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.AccountResponse, error) {
    account, err := h.service.CreateAccount(ctx, req.UserId, req.AccountType, req.Currency, req.InitialDeposit)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create account: %v", err)
    }

    return &pb.AccountResponse{
        Account: &pb.Account{
            Id:           account.ID,
            UserId:       account.UserID,
            AccountNumber: account.AccountNumber,
            AccountType:  account.AccountType,
            Balance:      account.Balance,
            Currency:     account.Currency,
            IsActive:     account.IsActive,
            CreatedAt:    account.CreatedAt.Format(time.RFC3339),
            UpdatedAt:    account.UpdatedAt.Format(time.RFC3339),
        },
    }, nil
}

// GetAccount implements the gRPC GetAccount method
func (h *AccountHandler) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.AccountResponse, error) {
    account, err := h.service.GetAccount(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "account not found: %v", err)
    }

    return &pb.AccountResponse{
        Account: &pb.Account{
            Id:           account.ID,
            UserId:       account.UserID,
            AccountNumber: account.AccountNumber,
            AccountType:  account.AccountType,
            Balance:      account.Balance,
            Currency:     account.Currency,
            IsActive:     account.IsActive,
            CreatedAt:    account.CreatedAt.Format(time.RFC3339),
            UpdatedAt:    account.UpdatedAt.Format(time.RFC3339),
        },
    }, nil
}

// And similar implementations for other gRPC methods
*/
