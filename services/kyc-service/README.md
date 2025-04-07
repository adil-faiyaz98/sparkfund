# KYC Service with AI Integration

The KYC (Know Your Customer) Service is a microservice responsible for managing customer identity verification and document processing in the SparkFund platform. It handles document uploads, verification processes, and maintains KYC profiles for users. The service now includes AI integration for enhanced verification capabilities.

## Features

- Document Management

  - Upload and store KYC documents
  - Support for multiple document types (ID, proof of address, etc.)
  - Document status tracking and expiration management
  - Secure document storage with encryption

- KYC Profile Management

  - Create and update customer profiles
  - Risk level assessment
  - Profile status tracking
  - Comprehensive customer information storage

- Verification Process

  - Document verification workflow
  - Multiple verification methods support
  - Verification history tracking
  - Confidence scoring system

- AI Integration

  - Document verification using AI models
  - Facial recognition for identity verification
  - Risk analysis based on user data and device information
  - Anomaly detection to identify suspicious behavior
  - AI model management and versioning

- Integration
  - AML service integration for enhanced screening
  - Notification service integration for status updates
  - User service integration for profile management

## Architecture

The service follows a clean architecture pattern with the following layers:

- **Domain Layer**: Contains business entities and logic
- **Service Layer**: Implements business use cases
- **Repository Layer**: Handles data persistence
- **Handler Layer**: Manages HTTP requests and responses
- **Infrastructure Layer**: Provides external dependencies and configurations

### Model Structure

The KYC service uses a clear separation between domain models and database models:

- **Domain Models**: Located in `internal/domain/`, these models represent business entities and contain business logic. They are free from ORM dependencies.
- **Database Models**: Located in `internal/model/`, these models are used for database operations and contain GORM tags.
- **Mappers**: Located in `internal/mapper/`, these functions convert between domain and database models.

For more information on the model structure, see the following documentation:

- [Model README](internal/model/README.md)
- [Domain README](internal/domain/README.md)
- [Mapper README](internal/mapper/README.md)
- [Migration Guide](internal/MIGRATION.md)

## API Endpoints

### Documents

- `POST /api/v1/documents` - Upload a new document
- `GET /api/v1/documents/:id` - Get document details
- `GET /api/v1/documents` - List documents
- `PUT /api/v1/documents/:id/status` - Update document status
- `DELETE /api/v1/documents/:id` - Delete document
- `GET /api/v1/documents/pending` - List pending documents
- `GET /api/v1/documents/expired` - List expired documents

### Profiles

- `POST /api/v1/profiles` - Create a new KYC profile
- `GET /api/v1/profiles/:user_id` - Get profile details
- `PUT /api/v1/profiles/:user_id` - Update profile
- `PUT /api/v1/profiles/:user_id/status` - Update profile status
- `PUT /api/v1/profiles/:user_id/risk-level` - Update risk level
- `GET /api/v1/profiles` - List profiles
- `GET /api/v1/profiles/stats` - Get profile statistics
- `DELETE /api/v1/profiles/:user_id` - Delete profile

### Verifications

- `POST /api/v1/verifications` - Create verification record
- `GET /api/v1/verifications/:id` - Get verification details
- `GET /api/v1/verifications/document/:document_id` - List document verifications
- `PUT /api/v1/verifications/:id` - Update verification
- `DELETE /api/v1/verifications/:id` - Delete verification
- `GET /api/v1/verifications/stats` - Get verification statistics
- `GET /api/v1/verifications/document/:document_id/history` - Get verification history
- `GET /api/v1/verifications/document/:document_id/summary` - Get verification summary

### AI Integration

- `GET /api/v1/ai/models` - List available AI models
- `POST /api/v1/ai/analyze-document` - Analyze a document using AI
- `POST /api/v1/ai/match-faces` - Match a selfie with a document photo
- `POST /api/v1/ai/analyze-risk` - Analyze risk based on user data
- `POST /api/v1/ai/detect-anomalies` - Detect anomalies in user behavior
- `POST /api/v1/ai/process-document` - Process a document through all AI checks

## Prerequisites

- Go 1.21 or later
- PostgreSQL 14 or later
- Redis (optional, for caching)
- Docker (optional, for containerization)

## Configuration

The service can be configured using environment variables or a `.env` file. See `.env.example` for available configuration options.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/sparkfund/kyc-service.git
cd kyc-service
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the service:

```bash
go run cmd/server/main.go
```

## Deployment

### Docker Deployment

#### Using Docker Compose (Recommended for Development)

1. Use Docker Compose to run the KYC service with PostgreSQL and all dependencies:

```bash
# On Linux/macOS
./run.sh

# On Windows
run.bat
```

This will start the KYC service with AI integration and a PostgreSQL database.

#### Manual Docker Deployment

1. Build the Docker image:

```bash
docker build -t sparkfund/kyc-service .
```

2. Run the container:

```bash
docker run -p 8080:8080 \
  --env-file .env \
  -v /path/to/documents:/data/documents \
  -v /path/to/models:/app/models \
  sparkfund/kyc-service
```

### Kubernetes Deployment (Recommended for Production)

The KYC service includes Kubernetes manifests for production deployment in the `k8s` directory:

- `deployment.yaml`: Defines the deployment with replicas, resource limits, and probes
- `service.yaml`: Exposes the service within the cluster
- `configmap.yaml`: Contains configuration for the service
- `secret.yaml`: Contains sensitive configuration (credentials, keys, etc.)
- `hpa.yaml`: Horizontal Pod Autoscaler for automatic scaling
- `pdb.yaml`: Pod Disruption Budget for high availability
- `network-policy.yaml`: Network security policies
- `service-monitor.yaml`: Prometheus ServiceMonitor for metrics collection
- `prometheus-rules.yaml`: Alerting rules for Prometheus

To deploy to Kubernetes:

```bash
# Apply all manifests
kubectl apply -f services/kyc-service/k8s/

# Check deployment status
kubectl -n sparkfund get pods -l app=kyc-service
```

#### Blue-Green Deployment

The service supports blue-green deployment for zero-downtime updates:

```bash
# Deploy to the inactive environment (blue or green)
kubectl -n sparkfund set image deployment/kyc-service-blue kyc-service=sparkfund/kyc-service:new-version

# Wait for deployment to complete
kubectl -n sparkfund rollout status deployment/kyc-service-blue

# Switch traffic to the new deployment
kubectl -n sparkfund patch service kyc-service -p '{"spec":{"selector":{"deployment":"blue"}}}'
```

## Development

### Running Tests

```bash
go test -v ./...
```

### Code Style

The project follows the standard Go code style. Run the linter:

```bash
golangci-lint run
```

### API Documentation

The API documentation is available in Swagger format at `/swagger-ui/` when running the service.

### Continuous Integration and Deployment

The KYC service uses GitHub Actions for CI/CD. The workflow is defined in `.github/workflows/kyc-service-ci-cd.yml` and includes:

#### CI Pipeline

1. **Linting**: Uses golangci-lint to check code quality
2. **Testing**: Runs unit tests and generates coverage reports
3. **Security Scanning**: Uses Gosec and Nancy to check for vulnerabilities
4. **SonarCloud Analysis**: Performs static code analysis
5. **Build**: Builds the Docker image and pushes it to the registry

#### CD Pipeline

1. **Staging Deployment**: Automatically deploys to the staging environment
2. **Integration Tests**: Runs integration tests against the staging environment
3. **Production Deployment**: Deploys to production using blue-green deployment
4. **Verification**: Verifies the deployment is healthy

The CI/CD pipeline includes several quality gates:

- Code coverage must be above 80%
- No critical or high security vulnerabilities
- All tests must pass
- SonarCloud quality gate must pass

## Security

- All endpoints require authentication with JWT
- Documents are encrypted at rest
- File uploads are validated and sanitized
- Rate limiting is implemented with configurable thresholds
- CORS is configured for secure cross-origin requests
- AI models are secured and access-controlled
- Device information is collected for risk analysis
- Anomaly detection identifies suspicious behavior
- Multi-factor authentication (MFA) support (configurable)
- Comprehensive logging for audit trails
- CSRF protection
- TLS/HTTPS support (configurable)
- Security headers (Content-Security-Policy, X-Frame-Options, etc.)
- Input validation for all parameters
- Database constraints and indexes for data integrity

## Monitoring and Observability

- Prometheus metrics exposed at `/metrics`
- Health checks at `/health`
- Readiness and liveness probes for Kubernetes
- Structured logging with configurable levels
- Distributed tracing with Jaeger
- Circuit breaker metrics for service resilience
- Alerting rules for Prometheus
- Integration with PagerDuty for production alerts
- Slack notifications for non-critical alerts
- Performance metrics collection
- Custom dashboards for Grafana

## Resilience and Performance

### Resilience

- **Circuit Breaker Pattern**: Prevents cascading failures by failing fast when dependencies are unavailable
- **Retry Mechanisms**: Automatically retries failed operations with exponential backoff
- **Fallback Strategies**: Provides alternative responses when primary operations fail
- **Timeout Management**: Configurable timeouts for all external calls
- **Graceful Degradation**: Continues to operate with reduced functionality when dependencies fail
- **Chaos Testing**: Regular chaos tests to verify resilience
- **Pod Anti-Affinity**: Ensures pods are distributed across nodes for high availability
- **Pod Disruption Budget**: Limits the number of pods that can be down simultaneously

### Performance

- **Connection Pooling**: Database connection pooling for efficient resource usage
- **Caching**: In-memory caching for frequently accessed data
- **Pagination**: All list endpoints support pagination
- **Efficient Database Queries**: Optimized queries with proper indexes
- **Resource Limits**: CPU and memory limits to prevent resource exhaustion
- **Horizontal Scaling**: Automatic scaling based on CPU and memory usage
- **Load Testing**: Regular load testing to identify bottlenecks
- **Performance Monitoring**: Continuous monitoring of response times and throughput

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
