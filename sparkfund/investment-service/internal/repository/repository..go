package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sparkfund/investment-service/internal/errors"
	"github.com/sparkfund/investment-service/internal/models"
	"go.uber.org/zap"
)

type Repository interface {
	GetClientInvestments(ctx context.Context, clientId string) ([]models.Investment, error)
	CreatePortfolio(ctx context.Context, portfolio models.Portfolio) (models.Portfolio, error)
	GetPortfolio(ctx context.Context, portfolioId string) (models.Portfolio, error)
	UpdatePortfolio(ctx context.Context, portfolioId string, portfolio models.Portfolio) (models.Portfolio, error)
	DeletePortfolio(ctx context.Context, portfolioId string) error
	GetClientWill(ctx context.Context, clientId string) (models.Will, error)
	CreateOrUpdateClientWill(ctx context.Context, clientId string, will models.Will) (models.Will, error)
	GetClientsWithdrawalThreshold(ctx context.Context, clientId string) (models.WithdrawalThreshold, error)
	SetClientWithdrawalThreshold(ctx context.Context, clientId string, threshold models.WithdrawalThreshold) (models.WithdrawalThreshold, error)
	UpdateClientWithdrawalThreshold(ctx context.Context, clientId string, threshold models.WithdrawalThreshold) (models.WithdrawalThreshold, error)
	GetInvestment(ctx context.Context, investmentId string) (models.Investment, error)
	UpdateInvestment(ctx context.Context, investmentId string, investment models.Investment) (models.Investment, error)
	DeleteInvestment(ctx context.Context, investmentId string) error
	GetInvestmentForecast(ctx context.Context, investmentId string) (models.InvestmentForecast, error)
	GetInvestmentRecommendation(ctx context.Context, investmentId string) (models.InvestmentRecommendation, error)
}

type PostgresRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPostgresRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &PostgresRepository{db: db, logger: logger}
}

func (r *PostgresRepository) GetClientInvestments(ctx context.Context, clientId string) ([]models.Investment, error) {
	query := `SELECT investment_id, client_id, portfolio_id, type, amount, purchase_date FROM investments WHERE client_id = $1`
	var investments []models.Investment
	err := r.db.SelectContext(ctx, &investments, query, clientId)
	if err != nil {
		r.logger.Error("Failed to get client investments", zap.Error(err), zap.String("clientId", clientId))
		return nil, errors.NewDatabaseError(err, errors.ComponentRepository)
	}
	return investments, nil
}

func (r *PostgresRepository) CreatePortfolio(ctx context.Context, portfolio models.Portfolio) (models.Portfolio, error) {
	query := `INSERT INTO portfolios (client_id, name, description) VALUES ($1, $2, $3) RETURNING portfolio_id, client_id, name, description`
	err := r.db.GetContext(ctx, &portfolio, query, portfolio.ClientId, portfolio.Name, portfolio.Description)
	if err != nil {
		r.logger.Error("Failed to create portfolio", zap.Error(err), zap.Int("clientId", portfolio.ClientId), zap.String("name", portfolio.Name))
		return models.Portfolio{}, errors.NewDatabaseError(err, errors.ComponentRepository)
	}
	return portfolio, nil
}

func (r *PostgresRepository) GetPortfolio(ctx context.Context, portfolioId string) (models.Portfolio, error) {
	query := `SELECT portfolio_id, client_id, name, description FROM portfolios WHERE portfolio_id = $1`
	var portfolio models.Portfolio
	err := r.db.GetContext(ctx, &portfolio, query, portfolioId)
	if err != nil {
		r.logger.Error("Failed to get portfolio", zap.Error(err), zap.String("portfolioId", portfolioId))
		return models.Portfolio{}, errors.NewDatabaseError(err, errors.ComponentRepository)
	}
	return portfolio, nil
}

func (r *PostgresRepository) UpdatePortfolio(ctx context.Context, portfolioId string, portfolio models.Portfolio) (models.Portfolio, error) {
	query := `UPDATE portfolios SET client_id = $1, name = $2, description = $3 WHERE portfolio_id = $4 RETURNING portfolio_id, client_id, name, description`
	err := r.db.GetContext(ctx, &portfolio, query, portfolio.ClientId, portfolio.Name, portfolio.Description, portfolioId)
	if err != nil {
		r.logger.Error("Failed to update portfolio", zap.Error(err), zap.String("portfolioId", portfolioId))
		return models.Portfolio{}, errors.NewDatabaseError(err, errors.ComponentRepository)
	}
	return portfolio, nil
}

func (r *PostgresRepository) DeletePortfolio(ctx context.Context, portfolioId string) error {
	query := `DELETE FROM portfolios WHERE portfolio_id = $1`
	_, err := r.db.ExecContext(ctx, query, portfolioId)
	if err != nil {
		r.logger.Error("Failed to delete portfolio", zap.Error(err), zap.String("portfolioId", portfolioId))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
	}
	return nil
}

func (r *PostgresRepository) GetClientWill(ctx context.Context, clientId string) (models.Will, error) {
	// Placeholder implementation
	r.logger.Info("Getting client will", zap.String("clientId", clientId))
	return models.Will{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) CreateOrUpdateClientWill(ctx context.Context, clientId string, will models.Will) (models.Will, error) {
	// Placeholder implementation
	r.logger.Info("Creating or updating client will", zap.String("clientId", clientId))
	return models.Will{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) GetClientsWithdrawalThreshold(ctx context.Context, clientId string) (models.WithdrawalThreshold, error) {
	// Placeholder implementation
	r.logger.Info("Getting client withdrawal threshold", zap.String("clientId", clientId))
	return models.WithdrawalThreshold{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) SetClientWithdrawalThreshold(ctx context.Context, clientId string, threshold models.WithdrawalThreshold) (models.WithdrawalThreshold, error) {
	// Placeholder implementation
	r.logger.Info("Setting client withdrawal threshold", zap.String("clientId", clientId))
	return models.WithdrawalThreshold{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) UpdateClientWithdrawalThreshold(ctx context.Context, clientId string, threshold models.WithdrawalThreshold) (models.WithdrawalThreshold, error) {
	// Placeholder implementation
	r.logger.Info("Updating client withdrawal threshold", zap.String("clientId", clientId))
	return models.WithdrawalThreshold{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) GetInvestment(ctx context.Context, investmentId string) (models.Investment, error) {
	// Placeholder implementation
	r.logger.Info("Getting investment", zap.String("investmentId", investmentId))
	return models.Investment{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) UpdateInvestment(ctx context.Context, investmentId string, investment models.Investment) (models.Investment, error) {
	// Placeholder implementation
	r.logger.Info("Updating investment", zap.String("investmentId", investmentId))
	return models.Investment{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) DeleteInvestment(ctx context.Context, investmentId string) error {
	// Placeholder implementation
	r.logger.Info("Deleting investment", zap.String("investmentId", investmentId))
	return fmt.Errorf("not implemented")
}

func (r *PostgresRepository) GetInvestmentForecast(ctx context.Context, investmentId string) (models.InvestmentForecast, error) {
	// Placeholder implementation
	r.logger.Info("Getting investment forecast", zap.String("investmentId", investmentId))
	return models.InvestmentForecast{}, fmt.Errorf("not implemented")
}

func (r *PostgresRepository) GetInvestmentRecommendation(ctx context.Context, investmentId string) (models.InvestmentRecommendation, error) {
	// Placeholder implementation
	r.logger.Info("Getting investment recommendation", zap.String("investmentId", investmentId))
	return models.InvestmentRecommendation{}, fmt.Errorf("not implemented")
}
