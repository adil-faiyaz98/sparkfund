package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *gorm.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

// Create creates a new audit event
func (r *PostgresRepository) Create(ctx context.Context, event *Event) error {
	// Convert request params to JSON
	var requestParamsJSON []byte
	if event.RequestParams != nil {
		var err error
		requestParamsJSON, err = json.Marshal(event.RequestParams)
		if err != nil {
			return fmt.Errorf("failed to marshal request params: %w", err)
		}
	}

	// Convert changes to JSON
	var changesJSON []byte
	if event.Changes != nil {
		var err error
		changesJSON, err = json.Marshal(event.Changes)
		if err != nil {
			return fmt.Errorf("failed to marshal changes: %w", err)
		}
	}

	// Convert metadata to JSON
	var metadataJSON []byte
	if event.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(event.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	// Create audit event in database
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO audit_events (
			id, timestamp, user_id, action, resource, resource_id, status,
			client_ip, user_agent, request_id, request_method, request_path,
			request_params, response_code, error_message, changes, metadata
		) VALUES (
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?, ?, ?
		)
	`,
		event.ID,
		event.Timestamp,
		event.UserID,
		event.Action,
		event.Resource,
		event.ResourceID,
		event.Status,
		event.ClientIP,
		event.UserAgent,
		event.RequestID,
		event.RequestMethod,
		event.RequestPath,
		requestParamsJSON,
		event.ResponseCode,
		event.ErrorMessage,
		changesJSON,
		metadataJSON,
	)

	if result.Error != nil {
		return fmt.Errorf("failed to create audit event: %w", result.Error)
	}

	return nil
}

// GetByID gets an audit event by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Event, error) {
	var event Event
	var requestParamsJSON, changesJSON, metadataJSON []byte

	result := r.db.WithContext(ctx).Raw(`
		SELECT
			id, timestamp, user_id, action, resource, resource_id, status,
			client_ip, user_agent, request_id, request_method, request_path,
			request_params, response_code, error_message, changes, metadata
		FROM audit_events
		WHERE id = ?
	`, id).Scan(&event.ID, &event.Timestamp, &event.UserID, &event.Action, &event.Resource,
		&event.ResourceID, &event.Status, &event.ClientIP, &event.UserAgent,
		&event.RequestID, &event.RequestMethod, &event.RequestPath, &requestParamsJSON,
		&event.ResponseCode, &event.ErrorMessage, &changesJSON, &metadataJSON)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get audit event: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	// Parse JSON fields
	if len(requestParamsJSON) > 0 {
		var requestParams interface{}
		if err := json.Unmarshal(requestParamsJSON, &requestParams); err != nil {
			return nil, fmt.Errorf("failed to unmarshal request params: %w", err)
		}
		event.RequestParams = requestParams
	}

	if len(changesJSON) > 0 {
		var changes interface{}
		if err := json.Unmarshal(changesJSON, &changes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal changes: %w", err)
		}
		event.Changes = changes
	}

	if len(metadataJSON) > 0 {
		var metadata interface{}
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		event.Metadata = metadata
	}

	return &event, nil
}

// GetByResourceID gets audit events by resource ID
func (r *PostgresRepository) GetByResourceID(ctx context.Context, resourceType ResourceType, resourceID string, limit, offset int) ([]*Event, int, error) {
	var events []*Event
	var count int64

	// Count total events
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM audit_events
		WHERE resource = ? AND resource_id = ?
	`, resourceType, resourceID).Scan(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit events: %w", err)
	}

	// Get events with pagination
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT
			id, timestamp, user_id, action, resource, resource_id, status,
			client_ip, user_agent, request_id, request_method, request_path,
			request_params, response_code, error_message, changes, metadata
		FROM audit_events
		WHERE resource = ? AND resource_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, resourceType, resourceID, limit, offset).Rows()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		var requestParamsJSON, changesJSON, metadataJSON []byte

		if err := rows.Scan(&event.ID, &event.Timestamp, &event.UserID, &event.Action, &event.Resource,
			&event.ResourceID, &event.Status, &event.ClientIP, &event.UserAgent,
			&event.RequestID, &event.RequestMethod, &event.RequestPath, &requestParamsJSON,
			&event.ResponseCode, &event.ErrorMessage, &changesJSON, &metadataJSON); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit event: %w", err)
		}

		// Parse JSON fields
		if len(requestParamsJSON) > 0 {
			var requestParams interface{}
			if err := json.Unmarshal(requestParamsJSON, &requestParams); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal request params: %w", err)
			}
			event.RequestParams = requestParams
		}

		if len(changesJSON) > 0 {
			var changes interface{}
			if err := json.Unmarshal(changesJSON, &changes); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal changes: %w", err)
			}
			event.Changes = changes
		}

		if len(metadataJSON) > 0 {
			var metadata interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			event.Metadata = metadata
		}

		events = append(events, &event)
	}

	return events, int(count), nil
}

// GetByUserID gets audit events by user ID
func (r *PostgresRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Event, int, error) {
	var events []*Event
	var count int64

	// Count total events
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(*)
		FROM audit_events
		WHERE user_id = ?
	`, userID).Scan(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit events: %w", err)
	}

	// Get events with pagination
	rows, err := r.db.WithContext(ctx).Raw(`
		SELECT
			id, timestamp, user_id, action, resource, resource_id, status,
			client_ip, user_agent, request_id, request_method, request_path,
			request_params, response_code, error_message, changes, metadata
		FROM audit_events
		WHERE user_id = ?
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, userID, limit, offset).Rows()

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		var requestParamsJSON, changesJSON, metadataJSON []byte

		if err := rows.Scan(&event.ID, &event.Timestamp, &event.UserID, &event.Action, &event.Resource,
			&event.ResourceID, &event.Status, &event.ClientIP, &event.UserAgent,
			&event.RequestID, &event.RequestMethod, &event.RequestPath, &requestParamsJSON,
			&event.ResponseCode, &event.ErrorMessage, &changesJSON, &metadataJSON); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit event: %w", err)
		}

		// Parse JSON fields
		if len(requestParamsJSON) > 0 {
			var requestParams interface{}
			if err := json.Unmarshal(requestParamsJSON, &requestParams); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal request params: %w", err)
			}
			event.RequestParams = requestParams
		}

		if len(changesJSON) > 0 {
			var changes interface{}
			if err := json.Unmarshal(changesJSON, &changes); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal changes: %w", err)
			}
			event.Changes = changes
		}

		if len(metadataJSON) > 0 {
			var metadata interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			event.Metadata = metadata
		}

		events = append(events, &event)
	}

	return events, int(count), nil
}

// Search searches for audit events
func (r *PostgresRepository) Search(ctx context.Context, query map[string]interface{}, limit, offset int) ([]*Event, int, error) {
	var events []*Event
	var count int64

	// Build WHERE clause
	var whereClause strings.Builder
	var args []interface{}

	if len(query) > 0 {
		whereClause.WriteString("WHERE ")
		first := true

		for key, value := range query {
			if !first {
				whereClause.WriteString(" AND ")
			}
			first = false

			switch key {
			case "start_time":
				whereClause.WriteString("timestamp >= ?")
				args = append(args, value)
			case "end_time":
				whereClause.WriteString("timestamp <= ?")
				args = append(args, value)
			case "action":
				whereClause.WriteString("action = ?")
				args = append(args, value)
			case "resource":
				whereClause.WriteString("resource = ?")
				args = append(args, value)
			case "status":
				whereClause.WriteString("status = ?")
				args = append(args, value)
			case "user_id":
				whereClause.WriteString("user_id = ?")
				args = append(args, value)
			case "resource_id":
				whereClause.WriteString("resource_id = ?")
				args = append(args, value)
			case "client_ip":
				whereClause.WriteString("client_ip = ?")
				args = append(args, value)
			case "request_id":
				whereClause.WriteString("request_id = ?")
				args = append(args, value)
			case "request_method":
				whereClause.WriteString("request_method = ?")
				args = append(args, value)
			case "request_path":
				whereClause.WriteString("request_path LIKE ?")
				args = append(args, "%"+value.(string)+"%")
			case "response_code":
				whereClause.WriteString("response_code = ?")
				args = append(args, value)
			}
		}
	}

	// Count total events
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM audit_events %s", whereClause.String())
	if err := r.db.WithContext(ctx).Raw(countQuery, args...).Scan(&count).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count audit events: %w", err)
	}

	// Get events with pagination
	selectQuery := fmt.Sprintf(`
		SELECT
			id, timestamp, user_id, action, resource, resource_id, status,
			client_ip, user_agent, request_id, request_method, request_path,
			request_params, response_code, error_message, changes, metadata
		FROM audit_events
		%s
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, whereClause.String())

	// Add limit and offset to args
	args = append(args, limit, offset)

	rows, err := r.db.WithContext(ctx).Raw(selectQuery, args...).Rows()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit events: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var event Event
		var requestParamsJSON, changesJSON, metadataJSON []byte

		if err := rows.Scan(&event.ID, &event.Timestamp, &event.UserID, &event.Action, &event.Resource,
			&event.ResourceID, &event.Status, &event.ClientIP, &event.UserAgent,
			&event.RequestID, &event.RequestMethod, &event.RequestPath, &requestParamsJSON,
			&event.ResponseCode, &event.ErrorMessage, &changesJSON, &metadataJSON); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit event: %w", err)
		}

		// Parse JSON fields
		if len(requestParamsJSON) > 0 {
			var requestParams interface{}
			if err := json.Unmarshal(requestParamsJSON, &requestParams); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal request params: %w", err)
			}
			event.RequestParams = requestParams
		}

		if len(changesJSON) > 0 {
			var changes interface{}
			if err := json.Unmarshal(changesJSON, &changes); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal changes: %w", err)
			}
			event.Changes = changes
		}

		if len(metadataJSON) > 0 {
			var metadata interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			event.Metadata = metadata
		}

		events = append(events, &event)
	}

	return events, int(count), nil
}
