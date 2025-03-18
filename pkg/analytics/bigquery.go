package analytics

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
)

// Client represents a BigQuery client
type Client struct {
	client    *bigquery.Client
	projectID string
	datasetID string
}

// Config holds BigQuery configuration
type Config struct {
	ProjectID string
	DatasetID string
	Location  string
}

// NewClient creates a new BigQuery client
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	client, err := bigquery.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create BigQuery client: %w", err)
	}

	return &Client{
		client:    client,
		projectID: cfg.ProjectID,
		datasetID: cfg.DatasetID,
	}, nil
}

// Transaction represents a financial transaction
type Transaction struct {
	UserID      string    `bigquery:"user_id"`
	Amount      float64   `bigquery:"amount"`
	Currency    string    `bigquery:"currency"`
	Category    string    `bigquery:"category"`
	Description string    `bigquery:"description"`
	Timestamp   time.Time `bigquery:"timestamp"`
}

// InsertTransactions inserts multiple transactions into BigQuery
func (c *Client) InsertTransactions(ctx context.Context, transactions []Transaction) error {
	inserter := c.client.Dataset(c.datasetID).Table("transactions").Inserter()
	return inserter.Put(ctx, transactions)
}

// DetectAnomalies detects anomalous transactions based on user's spending patterns
func (c *Client) DetectAnomalies(ctx context.Context, userID string, timeWindow time.Duration) ([]Transaction, error) {
	query := c.client.Query(`
		WITH UserStats AS (
			SELECT
				AVG(amount) as avg_amount,
				STDDEV(amount) as stddev_amount
			FROM ` + "`" + c.projectID + "." + c.datasetID + ".transactions`" + `
			WHERE user_id = @userID
			AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL @days DAY)
		)
		SELECT
			t.user_id,
			t.amount,
			t.currency,
			t.category,
			t.description,
			t.timestamp
		FROM ` + "`" + c.projectID + "." + c.datasetID + ".transactions`" + ` t
		CROSS JOIN UserStats
		WHERE t.user_id = @userID
		AND t.timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL @days DAY)
		AND ABS(t.amount - UserStats.avg_amount) > 2 * UserStats.stddev_amount
		ORDER BY t.timestamp DESC
	`)

	days := int64(timeWindow.Hours() / 24)
	query.Parameters = []bigquery.QueryParameter{
		{Name: "userID", Value: userID},
		{Name: "days", Value: days},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var anomalies []Transaction
	for {
		var t Transaction
		err := it.Next(&t)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}
		anomalies = append(anomalies, t)
	}

	return anomalies, nil
}

// GetSpendingTrends analyzes spending trends by category
func (c *Client) GetSpendingTrends(ctx context.Context, userID string, timeWindow time.Duration) (map[string]float64, error) {
	query := c.client.Query(`
		SELECT
			category,
			SUM(amount) as total_amount
		FROM ` + "`" + c.projectID + "." + c.datasetID + ".transactions`" + `
		WHERE user_id = @userID
		AND timestamp >= TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL @days DAY)
		GROUP BY category
		ORDER BY total_amount DESC
	`)

	days := int64(timeWindow.Hours() / 24)
	query.Parameters = []bigquery.QueryParameter{
		{Name: "userID", Value: userID},
		{Name: "days", Value: days},
	}

	it, err := query.Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	trends := make(map[string]float64)
	for {
		var row struct {
			Category    string
			TotalAmount float64
		}
		err := it.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}
		trends[row.Category] = row.TotalAmount
	}

	return trends, nil
}

// Close closes the BigQuery client
func (c *Client) Close() error {
	return c.client.Close()
}

// Example usage:
// client, err := analytics.NewClient(ctx, analytics.Config{
//     ProjectID: "my-project",
//     DatasetID: "financial_data",
// })
// if err != nil {
//     log.Fatal(err)
// }
// defer client.Close()
//
// transactions := []analytics.Transaction{...}
// err = client.InsertTransactions(ctx, transactions)
//
// anomalies, err := client.DetectAnomalies(ctx, userID, 30*24*time.Hour)
// trends, err := client.GetSpendingTrends(ctx, userID, 30*24*time.Hour)
