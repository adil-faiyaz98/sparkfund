module investment-service

go 1.21

require (
	github.gin-gonic/gin v1.9.1
	github.com/jmoiron/sqlx v1.3.5
	github.lib/pq v1.10.9
	github.prometheus/client_golang v1.18.0
	github.swaggo/files v1.0.0
	github.swaggo/gin-swagger v1.6.0
	github.swaggo/swag v0.0.0-20240129131215-5c267c66c426
	go.uber.org/zap v1.26.0
)

require (
	github.bytedance/sonic v1.9.1 // indirect
	github.cespare/xxhash/v2 v2.1.2 // indirect
	github.chenzhuoyu/base64x v0.0.0-20230716121745-296434df7cb3 // indirect
	github.chenzhuoyu/iasm v0.9.0 // indirect
	github.davecgh/go-spew v1.1.1 // indirect
	github.fsnotify/fsnotify v1.7.0 // indirect
	github.gabriel-vasile/mimetype v1.4.2 // indirect
	github.gin-contrib/sse v0.1.0 // indirect
	github.go-openapi/jsonpointer v0.20.2 // indirect
	github.go-openapi/jsonreference v0.20.4 // indirect
	github.go-openapi/spec v0.20.11 // indirect
	github.go-openapi/swag v0.22.9 // indirect
	github.go-playground/locales v0.14.1 // indirect
	github.go-playground/universal-translator v0.18.1 // indirect
	github.go-playground/validator/v10 v10.14.0 // indirect
	github.goccy/go-json v0.10.2 // indirect
	github.golang/protobuf v1.5.3 // indirect
	github.joho/godotenv v1.5.1 // indirect
	github.json-iterator/go v1.1.12 // indirect
	github.klauspost/cpuid/v2 v2.2.4 // indirect
	github.leodido/go-urn v1.2.4 // indirect
	github.mailru/easyjson v0.7.7 // indirect
	github.mattn/go-isatty v0.0.19 // indirect
	github.modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.modern-go/reflect2 v1.0.2 // indirect
	github.pelletier/go-toml/v2 v2.0.8 // indirect
	github.pmezard/go-difflib v1.0.0 // indirect
	github.stretchr/objx v0.5.0 // indirect
	github.swaggo/files/v2 v2.0.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// Add this to your go.mod file
replace github.KyleBanks/depth => github.com/KyleBanks/depth v1.2.1