package swagger

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed templates/*
var templates embed.FS

// SwaggerUI serves the Swagger UI interface
func SwaggerUI(openAPIPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFS(templates, "templates/swagger.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			OpenAPIPath string
		}{
			OpenAPIPath: openAPIPath,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// GenerateOpenAPIDoc generates OpenAPI documentation from Go code
func GenerateOpenAPIDoc(serviceName string) error {
	// This is a placeholder for actual OpenAPI generation
	// In a real implementation, you would use tools like swaggo/swag
	// to generate OpenAPI documentation from Go code comments

	// Create the OpenAPI directory if it doesn't exist
	apiDir := filepath.Join("api", serviceName, "v1")
	if err := os.MkdirAll(apiDir, 0755); err != nil {
		return err
	}

	return nil
}
