package docs

import "github.com/swaggo/swag"

var doc = `{
    "swagger": "3.0.0",
    "info": {
        "title": "Notification Service API",
        "description": "API for managing notifications, email, SMS, and push notifications",
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
            "url": "https://notifications.sparkfund.com",
            "description": "Production server"
        },
        {
            "url": "https://notifications.staging.sparkfund.com",
            "description": "Staging server"
        }
    ],
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
	Host:        "notifications.sparkfund.com",
	BasePath:    "/api/v1",
	Schemes:     []string{"https"},
	Title:       "Notification Service API",
	Description: "API for managing notifications, email, SMS, and push notifications",
}

func init() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		SwaggerTemplate: doc,
	})
} 