package docs

import "github.com/swaggo/swag"

var doc = `{
    "swagger": "3.0.0",
    "info": {
        "title": "AML Service API",
        "description": "API for managing Anti-Money Laundering (AML) screening and monitoring",
        "version": "1.0.0",
        "contact": {
            "name": "SparkFund API Team",
            "email": "api@sparkfund.com"
        },
        "license": {
            "name": "Proprietary",
            "url": "https://sparkfund.com/terms"
        }
    },
    "servers": [
        {
            "url": "https://aml.sparkfund.com",
            "description": "Production server"
        },
        {
            "url": "https://aml.staging.sparkfund.com",
            "description": "Staging server"
        }
    ],
    "paths": {
        "/api/v1/screening": {
            "post": {
                "summary": "Perform AML screening",
                "description": "Initiate AML screening for a user",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "requestBody": {
                    "required": true,
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/ScreeningRequest"
                            }
                        }
                    }
                },
                "responses": {
                    "201": {
                        "description": "Screening initiated successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/ScreeningResponse"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/Error"
                                }
                            }
                        }
                    }
                }
            },
            "get": {
                "summary": "List screenings",
                "description": "Retrieve a list of AML screenings",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "parameters": [
                    {
                        "in": "query",
                        "name": "status",
                        "schema": {
                            "type": "string",
                            "enum": ["pending", "completed", "failed"]
                        },
                        "description": "Filter by status"
                    },
                    {
                        "in": "query",
                        "name": "risk_level",
                        "schema": {
                            "type": "string",
                            "enum": ["low", "medium", "high"]
                        },
                        "description": "Filter by risk level"
                    },
                    {
                        "in": "query",
                        "name": "page",
                        "schema": {
                            "type": "integer",
                            "minimum": 1
                        },
                        "description": "Page number"
                    },
                    {
                        "in": "query",
                        "name": "page_size",
                        "schema": {
                            "type": "integer",
                            "minimum": 1,
                            "maximum": 100
                        },
                        "description": "Number of items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of screenings",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "type": "object",
                                    "properties": {
                                        "total": {
                                            "type": "integer"
                                        },
                                        "items": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/components/schemas/ScreeningResponse"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v1/screening/{id}": {
            "parameters": [
                {
                    "in": "path",
                    "name": "id",
                    "required": true,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    },
                    "description": "Screening ID"
                }
            ],
            "get": {
                "summary": "Get screening details",
                "description": "Retrieve detailed information about a specific screening",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Screening details",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/ScreeningResponse"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Screening not found"
                    }
                }
            }
        },
        "/api/v1/transactions": {
            "post": {
                "summary": "Process transaction",
                "description": "Process a transaction for AML checks",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "requestBody": {
                    "required": true,
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/Transaction"
                            }
                        }
                    }
                },
                "responses": {
                    "200": {
                        "description": "Transaction processed successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/Transaction"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Transaction blocked"
                    }
                }
            },
            "get": {
                "summary": "List transactions",
                "description": "Retrieve a list of transactions",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "parameters": [
                    {
                        "in": "query",
                        "name": "user_id",
                        "schema": {
                            "type": "string",
                            "format": "uuid"
                        },
                        "description": "Filter by user ID"
                    },
                    {
                        "in": "query",
                        "name": "status",
                        "schema": {
                            "type": "string",
                            "enum": ["pending", "completed", "failed", "blocked"]
                        },
                        "description": "Filter by status"
                    },
                    {
                        "in": "query",
                        "name": "start_date",
                        "schema": {
                            "type": "string",
                            "format": "date-time"
                        },
                        "description": "Start date for filtering"
                    },
                    {
                        "in": "query",
                        "name": "end_date",
                        "schema": {
                            "type": "string",
                            "format": "date-time"
                        },
                        "description": "End date for filtering"
                    },
                    {
                        "in": "query",
                        "name": "page",
                        "schema": {
                            "type": "integer",
                            "minimum": 1
                        },
                        "description": "Page number"
                    },
                    {
                        "in": "query",
                        "name": "page_size",
                        "schema": {
                            "type": "integer",
                            "minimum": 1,
                            "maximum": 100
                        },
                        "description": "Number of items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of transactions",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "type": "object",
                                    "properties": {
                                        "total": {
                                            "type": "integer"
                                        },
                                        "items": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/components/schemas/Transaction"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v1/risk-assessment/{user_id}": {
            "parameters": [
                {
                    "in": "path",
                    "name": "user_id",
                    "required": true,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    },
                    "description": "User ID"
                }
            ],
            "get": {
                "summary": "Get risk assessment",
                "description": "Retrieve risk assessment for a user",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Risk assessment details",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/RiskAssessment"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "User not found"
                    }
                }
            },
            "put": {
                "summary": "Update risk assessment",
                "description": "Update risk assessment for a user",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "requestBody": {
                    "required": true,
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/RiskAssessment"
                            }
                        }
                    }
                },
                "responses": {
                    "200": {
                        "description": "Risk assessment updated successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/RiskAssessment"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "User not found"
                    }
                }
            }
        },
        "/api/v1/alerts": {
            "get": {
                "summary": "List alerts",
                "description": "Retrieve a list of AML alerts",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "parameters": [
                    {
                        "in": "query",
                        "name": "type",
                        "schema": {
                            "type": "string",
                            "enum": ["suspicious_activity", "threshold_exceeded", "pattern_detected", "watchlist_match"]
                        },
                        "description": "Filter by alert type"
                    },
                    {
                        "in": "query",
                        "name": "severity",
                        "schema": {
                            "type": "string",
                            "enum": ["low", "medium", "high", "critical"]
                        },
                        "description": "Filter by severity"
                    },
                    {
                        "in": "query",
                        "name": "status",
                        "schema": {
                            "type": "string",
                            "enum": ["new", "in_review", "resolved", "dismissed"]
                        },
                        "description": "Filter by status"
                    },
                    {
                        "in": "query",
                        "name": "page",
                        "schema": {
                            "type": "integer",
                            "minimum": 1
                        },
                        "description": "Page number"
                    },
                    {
                        "in": "query",
                        "name": "page_size",
                        "schema": {
                            "type": "integer",
                            "minimum": 1,
                            "maximum": 100
                        },
                        "description": "Number of items per page"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of alerts",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "type": "object",
                                    "properties": {
                                        "total": {
                                            "type": "integer"
                                        },
                                        "items": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/components/schemas/Alert"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/api/v1/alerts/{id}": {
            "parameters": [
                {
                    "in": "path",
                    "name": "id",
                    "required": true,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    },
                    "description": "Alert ID"
                }
            ],
            "get": {
                "summary": "Get alert details",
                "description": "Retrieve detailed information about a specific alert",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Alert details",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/Alert"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Alert not found"
                    }
                }
            },
            "patch": {
                "summary": "Update alert status",
                "description": "Update the status of an alert",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "requestBody": {
                    "required": true,
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "required": ["status"],
                                "properties": {
                                    "status": {
                                        "type": "string",
                                        "enum": ["in_review", "resolved", "dismissed"]
                                    },
                                    "resolution": {
                                        "type": "object",
                                        "properties": {
                                            "action_taken": {
                                                "type": "string"
                                            },
                                            "notes": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "responses": {
                    "200": {
                        "description": "Alert updated successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/Alert"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid request"
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Alert not found"
                    }
                }
            }
        },
        "/api/v1/metrics": {
            "get": {
                "summary": "Get AML metrics",
                "description": "Retrieve metrics and statistics about AML operations",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "AML metrics",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/Metrics"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    }
                }
            }
        },
        "/health": {
            "get": {
                "summary": "Health check",
                "description": "Check the health status of the AML service",
                "responses": {
                    "200": {
                        "description": "Service is healthy",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "type": "object",
                                    "properties": {
                                        "status": {
                                            "type": "string",
                                            "enum": ["healthy", "degraded"]
                                        },
                                        "version": {
                                            "type": "string"
                                        },
                                        "uptime": {
                                            "type": "number"
                                        },
                                        "dependencies": {
                                            "type": "object",
                                            "properties": {
                                                "database": {
                                                    "type": "string",
                                                    "enum": ["healthy", "degraded", "down"]
                                                },
                                                "redis": {
                                                    "type": "string",
                                                    "enum": ["healthy", "degraded", "down"]
                                                },
                                                "watchlist_service": {
                                                    "type": "string",
                                                    "enum": ["healthy", "degraded", "down"]
                                                }
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    },
    "components": {
        "securitySchemes": {
            "BearerAuth": {
                "type": "http",
                "scheme": "bearer",
                "bearerFormat": "JWT"
            },
            "ApiKeyAuth": {
                "type": "apiKey",
                "in": "header",
                "name": "X-API-Key"
            }
        },
        "schemas": {
            "Error": {
                "type": "object",
                "properties": {
                    "code": {
                        "type": "string",
                        "description": "Error code"
                    },
                    "message": {
                        "type": "string",
                        "description": "Error message"
                    },
                    "details": {
                        "type": "object",
                        "description": "Additional error details"
                    }
                }
            },
            "ScreeningRequest": {
                "type": "object",
                "required": ["user_id", "screening_type", "customer_data"],
                "properties": {
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "ID of the user being screened"
                    },
                    "screening_type": {
                        "type": "string",
                        "enum": ["initial", "periodic", "enhanced"],
                        "description": "Type of screening to perform"
                    },
                    "customer_data": {
                        "type": "object",
                        "properties": {
                            "full_name": {
                                "type": "string",
                                "description": "Full legal name"
                            },
                            "date_of_birth": {
                                "type": "string",
                                "format": "date",
                                "description": "Date of birth"
                            },
                            "nationality": {
                                "type": "string",
                                "description": "Nationality"
                            },
                            "address": {
                                "type": "object",
                                "properties": {
                                    "street": {
                                        "type": "string"
                                    },
                                    "city": {
                                        "type": "string"
                                    },
                                    "state": {
                                        "type": "string"
                                    },
                                    "country": {
                                        "type": "string"
                                    },
                                    "postal_code": {
                                        "type": "string"
                                    }
                                }
                            },
                            "identification": {
                                "type": "object",
                                "properties": {
                                    "type": {
                                        "type": "string",
                                        "enum": ["passport", "national_id", "driving_license"]
                                    },
                                    "number": {
                                        "type": "string"
                                    },
                                    "issuing_country": {
                                        "type": "string"
                                    },
                                    "expiry_date": {
                                        "type": "string",
                                        "format": "date"
                                    }
                                }
                            },
                            "occupation": {
                                "type": "string",
                                "description": "Current occupation"
                            },
                            "source_of_funds": {
                                "type": "string",
                                "description": "Source of funds"
                            },
                            "expected_transaction_volume": {
                                "type": "number",
                                "description": "Expected monthly transaction volume"
                            },
                            "expected_transaction_frequency": {
                                "type": "string",
                                "enum": ["low", "medium", "high"],
                                "description": "Expected transaction frequency"
                            }
                        }
                    }
                }
            },
            "ScreeningResponse": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "Screening request ID"
                    },
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "ID of the screened user"
                    },
                    "status": {
                        "type": "string",
                        "enum": ["pending", "completed", "failed"],
                        "description": "Screening status"
                    },
                    "risk_level": {
                        "type": "string",
                        "enum": ["low", "medium", "high"],
                        "description": "Risk level assessment"
                    },
                    "risk_score": {
                        "type": "number",
                        "format": "float",
                        "minimum": 0,
                        "maximum": 100,
                        "description": "Numerical risk score"
                    },
                    "screening_results": {
                        "type": "object",
                        "properties": {
                            "watchlist_matches": {
                                "type": "array",
                                "items": {
                                    "type": "object",
                                    "properties": {
                                        "list_name": {
                                            "type": "string"
                                        },
                                        "match_type": {
                                            "type": "string"
                                        },
                                        "match_details": {
                                            "type": "string"
                                        }
                                    }
                                }
                            },
                            "pep_status": {
                                "type": "object",
                                "properties": {
                                    "is_pep": {
                                        "type": "boolean"
                                    },
                                    "pep_type": {
                                        "type": "string"
                                    },
                                    "jurisdiction": {
                                        "type": "string"
                                    }
                                }
                            },
                            "sanctions_matches": {
                                "type": "array",
                                "items": {
                                    "type": "object",
                                    "properties": {
                                        "sanction_list": {
                                            "type": "string"
                                        },
                                        "match_type": {
                                            "type": "string"
                                        },
                                        "match_details": {
                                            "type": "string"
                                        }
                                    }
                                }
                            }
                        }
                    },
                    "created_at": {
                        "type": "string",
                        "format": "date-time",
                        "description": "When screening was initiated"
                    },
                    "completed_at": {
                        "type": "string",
                        "format": "date-time",
                        "description": "When screening was completed"
                    },
                    "next_screening_date": {
                        "type": "string",
                        "format": "date-time",
                        "description": "When next screening is due"
                    }
                }
            },
            "Transaction": {
                "type": "object",
                "required": ["id", "user_id", "amount", "currency", "type"],
                "properties": {
                    "id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "Transaction ID"
                    },
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "User ID"
                    },
                    "amount": {
                        "type": "number",
                        "format": "float",
                        "description": "Transaction amount"
                    },
                    "currency": {
                        "type": "string",
                        "description": "Currency code"
                    },
                    "type": {
                        "type": "string",
                        "enum": ["deposit", "withdrawal", "transfer", "investment", "dividend"],
                        "description": "Transaction type"
                    },
                    "status": {
                        "type": "string",
                        "enum": ["pending", "completed", "failed", "blocked"],
                        "description": "Transaction status"
                    },
                    "timestamp": {
                        "type": "string",
                        "format": "date-time",
                        "description": "Transaction timestamp"
                    },
                    "source": {
                        "type": "object",
                        "properties": {
                            "account": {
                                "type": "string"
                            },
                            "institution": {
                                "type": "string"
                            }
                        }
                    },
                    "destination": {
                        "type": "object",
                        "properties": {
                            "account": {
                                "type": "string"
                            },
                            "institution": {
                                "type": "string"
                            }
                        }
                    },
                    "metadata": {
                        "type": "object",
                        "description": "Additional transaction metadata"
                    }
                }
            },
            "RiskAssessment": {
                "type": "object",
                "properties": {
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "User ID"
                    },
                    "risk_level": {
                        "type": "string",
                        "enum": ["low", "medium", "high"],
                        "description": "Current risk level"
                    },
                    "risk_score": {
                        "type": "number",
                        "format": "float",
                        "minimum": 0,
                        "maximum": 100,
                        "description": "Current risk score"
                    },
                    "factors": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "factor": {
                                    "type": "string",
                                    "description": "Risk factor"
                                },
                                "impact": {
                                    "type": "number",
                                    "format": "float",
                                    "description": "Impact on risk score"
                                },
                                "details": {
                                    "type": "string",
                                    "description": "Factor details"
                                }
                            }
                        }
                    },
                    "last_assessment": {
                        "type": "string",
                        "format": "date-time",
                        "description": "Last assessment timestamp"
                    },
                    "next_assessment": {
                        "type": "string",
                        "format": "date-time",
                        "description": "Next assessment due date"
                    },
                    "monitoring_level": {
                        "type": "string",
                        "enum": ["standard", "enhanced", "continuous"],
                        "description": "Current monitoring level"
                    }
                }
            },
            "Alert": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "Alert ID"
                    },
                    "type": {
                        "type": "string",
                        "enum": ["suspicious_activity", "threshold_exceeded", "pattern_detected", "watchlist_match"],
                        "description": "Alert type"
                    },
                    "severity": {
                        "type": "string",
                        "enum": ["low", "medium", "high", "critical"],
                        "description": "Alert severity"
                    },
                    "status": {
                        "type": "string",
                        "enum": ["new", "in_review", "resolved", "dismissed"],
                        "description": "Alert status"
                    },
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "Associated user ID"
                    },
                    "transaction_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "Associated transaction ID"
                    },
                    "details": {
                        "type": "object",
                        "description": "Alert details"
                    },
                    "created_at": {
                        "type": "string",
                        "format": "date-time",
                        "description": "When alert was created"
                    },
                    "updated_at": {
                        "type": "string",
                        "format": "date-time",
                        "description": "When alert was last updated"
                    },
                    "resolved_at": {
                        "type": "string",
                        "format": "date-time",
                        "description": "When alert was resolved"
                    },
                    "resolution": {
                        "type": "object",
                        "properties": {
                            "action_taken": {
                                "type": "string"
                            },
                            "notes": {
                                "type": "string"
                            },
                            "resolved_by": {
                                "type": "string",
                                "format": "uuid"
                            }
                        }
                    }
                }
            },
            "Metrics": {
                "type": "object",
                "properties": {
                    "total_screenings": {
                        "type": "integer",
                        "description": "Total number of screenings performed"
                    },
                    "screenings_by_status": {
                        "type": "object",
                        "properties": {
                            "completed": {
                                "type": "integer"
                            },
                            "pending": {
                                "type": "integer"
                            },
                            "failed": {
                                "type": "integer"
                            }
                        }
                    },
                    "risk_distribution": {
                        "type": "object",
                        "properties": {
                            "low": {
                                "type": "integer"
                            },
                            "medium": {
                                "type": "integer"
                            },
                            "high": {
                                "type": "integer"
                            }
                        }
                    },
                    "alerts_by_type": {
                        "type": "object",
                        "properties": {
                            "suspicious_activity": {
                                "type": "integer"
                            },
                            "threshold_exceeded": {
                                "type": "integer"
                            },
                            "pattern_detected": {
                                "type": "integer"
                            },
                            "watchlist_match": {
                                "type": "integer"
                            }
                        }
                    },
                    "average_processing_time": {
                        "type": "number",
                        "format": "float",
                        "description": "Average screening time in seconds"
                    }
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "aml.sparkfund.com",
	BasePath:    "/",
	Schemes:     []string{"https"},
	Title:       "AML Service API",
	Description: "API for managing Anti-Money Laundering (AML) screening and monitoring",
}

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  doc,
	})
}
