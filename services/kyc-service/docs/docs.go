// Package docs KYC Service API
//
// API for managing KYC verifications
//
//     Schemes: http
//     Host: localhost:8080
//     BasePath: /api/v1
//     Version: 1.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package docs

import (
	"embed"
)

//go:embed swagger.json
var SwaggerJSON embed.FS
