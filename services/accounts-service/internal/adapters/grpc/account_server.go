package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/adil-faiyaz98/money-pulse/proto/accounts/v1"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/domain"
	"github.com/adil-faiyaz98/money-pulse/services/accounts-service/internal/ports"
)

type Server struct {
	pb.UnimplementedAccountServiceServer
	api ports.APIPort
}

func NewServer(api ports.APIPort) *Server {
	return &Server{
		api: api,
	}
}

func (s *Server) Run(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAccountServiceServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func (s *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.Account, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	account, err := s.api.CreateAccount(ctx, userID, req.Name, domain.AccountType(req.Type), req.Currency)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(account), nil
}

func (s *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.Account, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account ID")
	}

	account, err := s.api.GetAccount(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(account), nil
}

func (s *Server) GetUserAccounts(ctx context.Context, req *pb.GetUserAccountsRequest) (*pb.GetUserAccountsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	accounts, err := s.api.GetUserAccounts(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoAccounts := make([]*pb.Account, len(accounts))
	for i, account := range accounts {
		protoAccounts[i] = toProto(account)
	}

	return &pb.GetUserAccountsResponse{
		Accounts: protoAccounts,
	}, nil
}

func (s *Server) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.Account, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account ID")
	}

	if err := s.api.UpdateAccount(ctx, id, req.Name, domain.AccountType(req.Type)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	account, err := s.api.GetAccount(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(account), nil
}

func (s *Server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account ID")
	}

	if err := s.api.DeleteAccount(ctx, id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteAccountResponse{}, nil
}

func (s *Server) GetAccountByNumber(ctx context.Context, req *pb.GetAccountByNumberRequest) (*pb.Account, error) {
	account, err := s.api.GetAccountByNumber(ctx, req.AccountNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(account), nil
}

func (s *Server) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*pb.Account, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account ID")
	}

	if err := s.api.UpdateBalance(ctx, id, req.Amount); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	account, err := s.api.GetAccount(ctx, id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return toProto(account), nil
}

func toProto(account *domain.Account) *pb.Account {
	return &pb.Account{
		Id:            account.ID.String(),
		UserId:        account.UserID.String(),
		Name:          account.Name,
		AccountNumber: account.AccountNumber,
		Type:          pb.AccountType(pb.AccountType_value[string(account.Type)]),
		Balance:       account.Balance,
		Currency:      account.Currency,
		CreatedAt:     account.CreatedAt.Unix(),
		UpdatedAt:     account.UpdatedAt.Unix(),
	}
}
