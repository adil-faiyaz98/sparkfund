module investment-service

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/go-chi/chi/v5 v5.0.12
	github.com/go-chi/cors v1.2.1
	github.com/google/uuid v1.6.0
	github.com/prometheus/client_golang v1.17.0
	github.com/sirupsen/logrus v1.9.3
	github.com/sony/gobreaker v1.0.0
	go.uber.org/zap v1.26.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)

// Pinned indirect dependencies for security
require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.9 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
)
