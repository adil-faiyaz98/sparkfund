package docs

import "github.com/swaggo/swag"

var doc = `{
    "swagger": "3.0.0",
    "info": {
        "title": "KYC Service API",
        "description": "API for managing Know Your Customer (KYC) verification and document processing",
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
            "url": "https://kyc.sparkfund.com",
            "description": "Production server"
        },
        {
            "url": "https://kyc.staging.sparkfund.com",
            "description": "Staging server"
        }
    ],
    "paths": {
        "/api/v1/documents": {
            "post": {
                "summary": "Upload document",
                "description": "Upload a new KYC document for verification",
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
                                "$ref": "#/components/schemas/DocumentUploadRequest"
                            }
                        }
                    }
                },
                "responses": {
                    "201": {
                        "description": "Document uploaded successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/DocumentResponse"
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
                "summary": "List documents",
                "description": "Retrieve a list of KYC documents",
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
                        "name": "document_type",
                        "schema": {
                            "type": "string",
                            "enum": ["passport", "national_id", "driving_license", "proof_of_address", "bank_statement", "tax_return", "employment_contract", "utility_bill"]
                        },
                        "description": "Filter by document type"
                    },
                    {
                        "in": "query",
                        "name": "status",
                        "schema": {
                            "type": "string",
                            "enum": ["pending", "verified", "rejected", "expired"]
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
                        "description": "List of documents",
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
                                                "$ref": "#/components/schemas/DocumentResponse"
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
        "/api/v1/documents/{id}": {
            "parameters": [
                {
                    "in": "path",
                    "name": "id",
                    "required": true,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    },
                    "description": "Document ID"
                }
            ],
            "get": {
                "summary": "Get document details",
                "description": "Retrieve detailed information about a specific document",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Document details",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/DocumentResponse"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Document not found"
                    }
                }
            },
            "patch": {
                "summary": "Update document status",
                "description": "Update the verification status of a document",
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
                                "$ref": "#/components/schemas/VerificationRequest"
                            }
                        }
                    }
                },
                "responses": {
                    "200": {
                        "description": "Document updated successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/DocumentResponse"
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
                        "description": "Document not found"
                    }
                }
            }
        },
        "/api/v1/profiles/{user_id}": {
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
                "summary": "Get KYC profile",
                "description": "Retrieve the complete KYC profile for a user",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "KYC profile details",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/KYCProfile"
                                }
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "404": {
                        "description": "Profile not found"
                    }
                }
            },
            "put": {
                "summary": "Update KYC profile",
                "description": "Update the KYC profile information",
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
                                "$ref": "#/components/schemas/KYCProfile"
                            }
                        }
                    }
                },
                "responses": {
                    "200": {
                        "description": "Profile updated successfully",
                        "content": {
                            "application/json": {
                                "schema": {
                                    "$ref": "#/components/schemas/KYCProfile"
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
                        "description": "Profile not found"
                    }
                }
            }
        },
        "/api/v1/metrics": {
            "get": {
                "summary": "Get KYC metrics",
                "description": "Retrieve metrics and statistics about KYC operations",
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "responses": {
                    "200": {
                        "description": "KYC metrics",
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
                "description": "Check the health status of the KYC service",
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
                                                "document_processing_service": {
                                                    "type": "string",
                                                    "enum": ["healthy", "degraded", "down"]
                                                },
                                                "ocr_service": {
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
            "DocumentUploadRequest": {
                "type": "object",
                "required": ["user_id", "document_type", "file_data"],
                "properties": {
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "ID of the user submitting the document"
                    },
                    "document_type": {
                        "type": "string",
                        "enum": ["passport", "national_id", "driving_license", "proof_of_address", "bank_statement", "tax_return", "employment_contract", "utility_bill"],
                        "description": "Type of document being uploaded"
                    },
                    "file_data": {
                        "type": "string",
                        "format": "base64",
                        "description": "Base64 encoded document file"
                    },
                    "metadata": {
                        "type": "object",
                        "properties": {
                            "issuing_country": {
                                "type": "string"
                            },
                            "expiry_date": {
                                "type": "string",
                                "format": "date"
                            },
                            "document_number": {
                                "type": "string"
                            },
                            "additional_info": {
                                "type": "object"
                            }
                        }
                    }
                }
            },
            "DocumentResponse": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "Document ID"
                    },
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "User ID"
                    },
                    "document_type": {
                        "type": "string",
                        "enum": ["passport", "national_id", "driving_license", "proof_of_address", "bank_statement", "tax_return", "employment_contract", "utility_bill"],
                        "description": "Type of document"
                    },
                    "status": {
                        "type": "string",
                        "enum": ["pending", "verified", "rejected", "expired"],
                        "description": "Document verification status"
                    },
                    "verification_details": {
                        "type": "object",
                        "properties": {
                            "verified_by": {
                                "type": "string",
                                "format": "uuid"
                            },
                            "verified_at": {
                                "type": "string",
                                "format": "date-time"
                            },
                            "verification_method": {
                                "type": "string",
                                "enum": ["manual", "automated", "third_party"]
                            },
                            "confidence_score": {
                                "type": "number",
                                "format": "float",
                                "minimum": 0,
                                "maximum": 100
                            },
                            "rejection_reason": {
                                "type": "string"
                            }
                        }
                    },
                    "metadata": {
                        "type": "object",
                        "properties": {
                            "issuing_country": {
                                "type": "string"
                            },
                            "expiry_date": {
                                "type": "string",
                                "format": "date"
                            },
                            "document_number": {
                                "type": "string"
                            },
                            "additional_info": {
                                "type": "object"
                            }
                        }
                    },
                    "created_at": {
                        "type": "string",
                        "format": "date-time"
                    },
                    "updated_at": {
                        "type": "string",
                        "format": "date-time"
                    }
                }
            },
            "KYCProfile": {
                "type": "object",
                "properties": {
                    "user_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "User ID"
                    },
                    "status": {
                        "type": "string",
                        "enum": ["pending", "in_review", "approved", "rejected", "suspended"],
                        "description": "Overall KYC status"
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
                    "personal_info": {
                        "type": "object",
                        "properties": {
                            "full_name": {
                                "type": "string"
                            },
                            "date_of_birth": {
                                "type": "string",
                                "format": "date"
                            },
                            "nationality": {
                                "type": "string"
                            },
                            "tax_id": {
                                "type": "string"
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
                            }
                        }
                    },
                    "employment_info": {
                        "type": "object",
                        "properties": {
                            "occupation": {
                                "type": "string"
                            },
                            "employer": {
                                "type": "string"
                            },
                            "employment_status": {
                                "type": "string",
                                "enum": ["employed", "self_employed", "retired", "unemployed"]
                            },
                            "annual_income": {
                                "type": "number",
                                "format": "float"
                            },
                            "source_of_funds": {
                                "type": "string"
                            }
                        }
                    },
                    "financial_info": {
                        "type": "object",
                        "properties": {
                            "expected_transaction_volume": {
                                "type": "number",
                                "format": "float"
                            },
                            "expected_transaction_frequency": {
                                "type": "string",
                                "enum": ["low", "medium", "high"]
                            },
                            "investment_experience": {
                                "type": "string",
                                "enum": ["beginner", "intermediate", "advanced"]
                            },
                            "investment_goals": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "documents": {
                        "type": "array",
                        "items": {
                            "$ref": "#/components/schemas/DocumentResponse"
                        }
                    },
                    "created_at": {
                        "type": "string",
                        "format": "date-time"
                    },
                    "updated_at": {
                        "type": "string",
                        "format": "date-time"
                    },
                    "last_review_date": {
                        "type": "string",
                        "format": "date-time"
                    },
                    "next_review_date": {
                        "type": "string",
                        "format": "date-time"
                    }
                }
            },
            "VerificationRequest": {
                "type": "object",
                "required": ["document_id", "action"],
                "properties": {
                    "document_id": {
                        "type": "string",
                        "format": "uuid",
                        "description": "ID of the document to verify"
                    },
                    "action": {
                        "type": "string",
                        "enum": ["approve", "reject", "request_additional_info"],
                        "description": "Verification action"
                    },
                    "notes": {
                        "type": "string",
                        "description": "Additional notes or comments"
                    },
                    "additional_info_required": {
                        "type": "object",
                        "description": "Details of additional information required"
                    }
                }
            },
            "Metrics": {
                "type": "object",
                "properties": {
                    "total_profiles": {
                        "type": "integer",
                        "description": "Total number of KYC profiles"
                    },
                    "profiles_by_status": {
                        "type": "object",
                        "properties": {
                            "pending": {
                                "type": "integer"
                            },
                            "in_review": {
                                "type": "integer"
                            },
                            "approved": {
                                "type": "integer"
                            },
                            "rejected": {
                                "type": "integer"
                            },
                            "suspended": {
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
                    "average_processing_time": {
                        "type": "number",
                        "format": "float",
                        "description": "Average document verification time in seconds"
                    },
                    "documents_by_type": {
                        "type": "object",
                        "properties": {
                            "passport": {
                                "type": "integer"
                            },
                            "national_id": {
                                "type": "integer"
                            },
                            "driving_license": {
                                "type": "integer"
                            },
                            "proof_of_address": {
                                "type": "integer"
                            },
                            "bank_statement": {
                                "type": "integer"
                            },
                            "tax_return": {
                                "type": "integer"
                            },
                            "employment_contract": {
                                "type": "integer"
                            },
                            "utility_bill": {
                                "type": "integer"
                            }
                        }
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
	Host:        "kyc.sparkfund.com",
	BasePath:    "/",
	Schemes:     []string{"https"},
	Title:       "KYC Service API",
	Description: "API for managing Know Your Customer (KYC) verification and document processing",
}

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate:  doc,
	})
}
