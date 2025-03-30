package repositories_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sparkfund/services/investment-service/internal/models"
	"github.com/sparkfund/services/investment-service/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	// Create a sqlmock database connection
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Create a gorm DB instance which uses the mock connection
	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	return gdb, mock, db
}

func TestCreateInvestment(t *testing.T) {
	// Setup
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repositories.NewInvestmentRepository(db)
	ctx := context.Background()

	// Create test investment
	now := time.Now()
	investment := &models.Investment{
		UserID:        1,
		PortfolioID:   1,
		Amount:        1000.0,
		Type:          "STOCK",
		Status:        "ACTIVE",
		PurchaseDate:  now,
		PurchasePrice: 150.0,
		Symbol:        "AAPL",
		Quantity:      10,
		Notes:         "Test investment",
	}

	// Set up mock expectations - matches the SQL query that GORM will execute
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "investments" ("user_id","portfolio_id","amount","type","status","purchase_date","purchase_price","symbol","quantity","notes","created_at","updated_at","sell_date","sell_price") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) RETURNING "id"`)).
		WithArgs(investment.UserID, investment.PortfolioID, investment.Amount, investment.Type, investment.Status, investment.PurchaseDate, investment.PurchasePrice, investment.Symbol, investment.Quantity, investment.Notes, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Call the repository method
	err := repo.Create(ctx, investment)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, uint(1), investment.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	// Setup
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repositories.NewInvestmentRepository(db)
	ctx := context.Background()

	// Set up test data
	now := time.Now()
	expectedInvestment := &models.Investment{
		ID:            1,
		UserID:        1,
		PortfolioID:   1,
		Amount:        1000.0,
		Type:          "STOCK",
		Status:        "ACTIVE",
		PurchaseDate:  now,
		PurchasePrice: 150.0,
		Symbol:        "AAPL",
		Quantity:      10,
		Notes:         "Test investment",
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Set up mock expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "portfolio_id", "amount", "type", "status", "purchase_date", "purchase_price", "symbol", "quantity", "notes", "created_at", "updated_at", "sell_date", "sell_price"}).
		AddRow(expectedInvestment.ID, expectedInvestment.UserID, expectedInvestment.PortfolioID, expectedInvestment.Amount, expectedInvestment.Type, expectedInvestment.Status, expectedInvestment.PurchaseDate, expectedInvestment.PurchasePrice, expectedInvestment.Symbol, expectedInvestment.Quantity, expectedInvestment.Notes, expectedInvestment.CreatedAt, expectedInvestment.UpdatedAt, nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "investments" WHERE id = $1 ORDER BY "investments"."id" LIMIT 1`)).
		WithArgs(1).
		WillReturnRows(rows)

	// Call the repository method
	investment, err := repo.GetByID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, investment)
	assert.Equal(t, expectedInvestment.ID, investment.ID)
	assert.Equal(t, expectedInvestment.Symbol, investment.Symbol)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByUserID(t *testing.T) {
	// Setup
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repositories.NewInvestmentRepository(db)
	ctx := context.Background()

	// Set up test data
	now := time.Now()
	expectedInvestments := []models.Investment{
		{
			ID:            1,
			UserID:        1,
			PortfolioID:   1,
			Amount:        1000.0,
			Type:          "STOCK",
			Status:        "ACTIVE",
			PurchaseDate:  now,
			PurchasePrice: 150.0,
			Symbol:        "AAPL",
			Quantity:      10,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            2,
			UserID:        1,
			PortfolioID:   1,
			Amount:        2000.0,
			Type:          "STOCK",
			Status:        "ACTIVE",
			PurchaseDate:  now,
			PurchasePrice: 200.0,
			Symbol:        "GOOGL",
			Quantity:      5,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	// Set up mock expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "portfolio_id", "amount", "type", "status", "purchase_date", "purchase_price", "symbol", "quantity", "notes", "created_at", "updated_at", "sell_date", "sell_price"})
	for _, inv := range expectedInvestments {
		rows.AddRow(inv.ID, inv.UserID, inv.PortfolioID, inv.Amount, inv.Type, inv.Status, inv.PurchaseDate, inv.PurchasePrice, inv.Symbol, inv.Quantity, inv.Notes, inv.CreatedAt, inv.UpdatedAt, nil, nil)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "investments" WHERE user_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	// Call the repository method
	investments, err := repo.GetByUserID(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, investments, 2)
	assert.Equal(t, expectedInvestments[0].ID, investments[0].ID)
	assert.Equal(t, expectedInvestments[1].ID, investments[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	// Setup
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repositories.NewInvestmentRepository(db)
	ctx := context.Background()

	// Create test investment
	now := time.Now()
	investment := &models.Investment{
		ID:            1,
		UserID:        1,
		PortfolioID:   1,
		Amount:        1500.0, // Updated amount
		Type:          "STOCK",
		Status:        "ACTIVE",
		PurchaseDate:  now,
		PurchasePrice: 150.0,
		Symbol:        "AAPL",
		Quantity:      10,
		Notes:         "Updated notes",
		UpdatedAt:     now,
	}

	// Set up mock expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "investments" SET "user_id"=$1,"portfolio_id"=$2,"amount"=$3,"type"=$4,"status"=$5,"purchase_date"=$6,"purchase_price"=$7,"symbol"=$8,"quantity"=$9,"notes"=$10,"created_at"=$11,"updated_at"=$12,"sell_date"=$13,"sell_price"=$14 WHERE "id" = $15`)).
		WithArgs(investment.UserID, investment.PortfolioID, investment.Amount, investment.Type, investment.Status, investment.PurchaseDate, investment.PurchasePrice, investment.Symbol, investment.Quantity, investment.Notes, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil, investment.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository method
	err := repo.Update(ctx, investment)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	// Setup
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repositories.NewInvestmentRepository(db)
	ctx := context.Background()

	// Set up mock expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "investments" WHERE "investments"."id" = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Call the repository method
	err := repo.Delete(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAll(t *testing.T) {
	// Setup
	db, mock, sqlDB := setupMockDB(t)
	defer sqlDB.Close()

	repo := repositories.NewInvestmentRepository(db)
	ctx := context.Background()

	// Set up test data
	now := time.Now()
	expectedInvestments := []models.Investment{
		{
			ID:            1,
			UserID:        1,
			PortfolioID:   1,
			Amount:        1000.0,
			Type:          "STOCK",
			Status:        "ACTIVE",
			PurchaseDate:  now,
			PurchasePrice: 150.0,
			Symbol:        "AAPL",
			Quantity:      10,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            2,
			UserID:        1,
			PortfolioID:   1,
			Amount:        2000.0,
			Type:          "STOCK",
			Status:        "ACTIVE",
			PurchaseDate:  now,
			PurchasePrice: 200.0,
			Symbol:        "GOOGL",
			Quantity:      5,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	// Set up mock expectations for count
	countRows := sqlmock.NewRows([]string{"count"}).AddRow(10) // Total 10 investments
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "investments"`)).
		WillReturnRows(countRows)

	// Set up mock expectations for pagination
	rows := sqlmock.NewRows([]string{"id", "user_id", "portfolio_id", "amount", "type", "status", "purchase_date", "purchase_price", "symbol", "quantity", "notes", "created_at", "updated_at", "sell_date", "sell_price"})
	for _, inv := range expectedInvestments {
		rows.AddRow(inv.ID, inv.UserID, inv.PortfolioID, inv.Amount, inv.Type, inv.Status, inv.PurchaseDate, inv.PurchasePrice, inv.Symbol, inv.Quantity, inv.Notes, inv.CreatedAt, inv.UpdatedAt, nil, nil)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "investments" LIMIT 10 OFFSET 0`)).
		WillReturnRows(rows)

	// Call the repository method
	investments, total, err := repo.GetAll(ctx, 1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, investments, 2)
	assert.Equal(t, int64(10), total)
	assert.Equal(t, expectedInvestments[0].ID, investments[0].ID)
	assert.Equal(t, expectedInvestments[1].ID, investments[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}
