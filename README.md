# SparkFund Platform

<p align="center">
  <img src="https://via.placeholder.com/200x200?text=SparkFund" alt="SparkFund Logo" width="200" height="200">
</p>

SparkFund is a comprehensive financial platform with AI-powered investment recommendations, KYC verification, and user management capabilities. Built with a microservices architecture, it provides a scalable and robust foundation for financial applications.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Services](#-services)
- [AI Capabilities](#-ai-capabilities)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [Cloud Deployment](#-cloud-deployment)
- [Infrastructure as Code](#-infrastructure-as-code)
- [Containerization & Orchestration](#-containerization--orchestration)
- [Cost Management](#-cost-management)
- [Testing](#-testing)
- [Monitoring and Observability](#-monitoring-and-observability)
- [Troubleshooting](#-troubleshooting)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- **Microservices Architecture**: Scalable, maintainable, and resilient system design
- **AI-Powered Investment Recommendations**: Advanced algorithms for personalized investment advice (POC)
- **KYC Verification**: Secure and compliant customer onboarding with AI-powered document verification (POC)
- **User Management**: Comprehensive user authentication and profile management
- **API Gateway**: Centralized entry point with routing, authentication, and rate limiting
- **Swagger Documentation**: Interactive API documentation for all services
- **Cloud-Native Design**: Ready for deployment on major cloud providers (AWS, Azure, GCP)
- **Infrastructure as Code**: Terraform configurations for automated infrastructure provisioning
- **Kubernetes Orchestration**: Helm charts for streamlined deployment and management
- **Cost Optimization**: Kubecost integration for Kubernetes cost monitoring and optimization
- **Monitoring and Observability**: Prometheus, Grafana, and Jaeger integration

## 🏗️ Architecture

SparkFund follows a microservices architecture with the following components:

<p align="center">
  <img src="https://via.placeholder.com/800x400?text=SparkFund+Architecture" alt="SparkFund Architecture" width="800">
</p>

1. **API Gateway**: Routes requests to appropriate services, handles authentication, and rate limiting
2. **KYC Service**: Handles customer verification and onboarding with AI-powered document verification
3. **Investment Service**: Manages investments and provides AI-powered recommendations
4. **User Service**: Manages user accounts, authentication, and profiles
5. **AI Service**: Provides AI capabilities for document verification, facial recognition, and investment analysis
6. **Supporting Infrastructure**: PostgreSQL, Redis, Prometheus, Grafana, and Jaeger

## 🚀 Services

| Service            | Description                                                                                                  |
| ------------------ | ------------------------------------------------------------------------------------------------------------ |
| API Gateway        | Routes requests to appropriate services, handles authentication, rate limiting, and request validation        |
| KYC Service        | Handles customer verification and onboarding with AI-powered document verification and identity validation    |
| Investment Service | Manages investments, portfolios, and provides AI-powered recommendations and market analysis                 |
| User Service       | Manages user accounts, authentication, profiles, and permissions                                             |
| AI Service         | Provides AI capabilities for document verification, facial recognition, NLP, and investment analysis         |

## 🧠 AI Capabilities

SparkFund integrates advanced AI capabilities across its services as a **Proof of Concept (POC)**. While the core microservices functionality is fully implemented, the AI features demonstrate how machine learning can enhance financial services.

### Investment Service AI (POC Implementation)

- **Recommendation System**: Personalized investment recommendations based on user profile and market conditions
  - *Models*: Collaborative filtering, content-based filtering, and hybrid recommendation systems
  - *Implementation*: API endpoints with sample responses that simulate model predictions
  - *Endpoints*: `/api/v1/ai/recommendations`, `/api/v1/ai/recommendations/personalized/:userId`

- **Natural Language Processing**: Analysis of market news for investment signals
  - *Models*: BERT-based sentiment analysis, named entity recognition, and topic modeling
  - *Implementation*: Text processing pipeline with pre-trained models
  - *Endpoints*: `/api/v1/ai/advanced/news/analyze`

- **Time Series Forecasting**: Price prediction and market movement analysis
  - *Models*: ARIMA, LSTM, and Prophet for time series forecasting
  - *Implementation*: Forecasting endpoints with sample historical data
  - *Endpoints*: `/api/v1/ai/advanced/price/forecast`, `/api/v1/ai/advanced/market/predict`

- **Reinforcement Learning**: Portfolio optimization and investment strategy learning
  - *Models*: Deep Q-Networks (DQN) and Proximal Policy Optimization (PPO)
  - *Implementation*: Simulated environment for portfolio optimization
  - *Endpoints*: `/api/v1/ai/advanced/portfolio/optimize`, `/api/v1/ai/advanced/portfolio/action`

- **Anomaly Detection**: Fraud detection and security monitoring
  - *Models*: Isolation Forest, One-Class SVM, and Autoencoder-based anomaly detection
  - *Implementation*: Transaction monitoring system with anomaly scoring
  - *Endpoints*: `/api/v1/ai/security/analyze`

### KYC Service AI (POC Implementation)

- **Document Verification**: AI-powered verification of identity documents
  - *Models*: YOLOv5 for object detection, OCR with Tesseract and PaddleOCR
  - *Implementation*: Document processing pipeline with validation rules
  - *Endpoints*: `/api/v1/verification/document`

- **Facial Recognition**: Face matching and liveness detection
  - *Models*: FaceNet for face embeddings, anti-spoofing models for liveness detection
  - *Implementation*: Face comparison and liveness check endpoints
  - *Endpoints*: `/api/v1/verification/facial`

- **Risk Analysis**: AI-powered risk assessment for KYC verification
  - *Models*: Gradient Boosting models (XGBoost, LightGBM) for risk scoring
  - *Implementation*: Risk assessment pipeline with multiple data sources
  - *Endpoints*: `/api/v1/verification/risk`

### AI Service (Central AI Capabilities)

- **Document Analysis**: Extract information from documents using computer vision and NLP
  - *Models*: CRAFT for text detection, BERT for information extraction
  - *Implementation*: End-to-end document processing pipeline
  - *Endpoints*: `/api/v1/document/analyze`

- **Image Classification**: Classify images for various purposes
  - *Models*: EfficientNet, ResNet for image classification
  - *Implementation*: Image processing and classification endpoints
  - *Endpoints*: `/api/v1/image/classify`

- **Text Analysis**: Sentiment analysis, entity extraction, and text classification
  - *Models*: DistilBERT for sentiment analysis, SpaCy for NER
  - *Implementation*: Text processing endpoints with multiple analysis options
  - *Endpoints*: `/api/v1/text/analyze`

### POC Implementation Notes

The AI capabilities are implemented as a proof of concept with the following characteristics:

1. **API Endpoints**: All AI endpoints are fully defined and documented in Swagger
2. **Sample Responses**: Endpoints return realistic sample responses based on the expected model output
3. **Model Placeholders**: Code structure is in place for integrating actual models
4. **Demonstration Ready**: The system demonstrates the potential of AI in financial services
5. **Production Path**: Clear path for replacing POC implementations with production models

## 🚦 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### Clone the Repository

```bash
git clone https://github.com/yourusername/sparkfund.git
cd sparkfund
```

### Starting All Services

```bash
# Start all services
docker-compose up -d

# Alternatively, use the provided scripts
# Windows
.\scripts\run-all-services.bat

# Linux/Mac
./scripts/run-all-services.sh
```

This will start all services and supporting infrastructure.

### Service URLs

Once all services are running, you can access them at:

| Service            | URL                    | Description                                   |
| ------------------ | ---------------------- | --------------------------------------------- |
| API Gateway        | http://localhost:8080  | Main entry point for all services             |
| KYC Service        | http://localhost:8081  | Know Your Customer verification               |
| Investment Service | http://localhost:8082  | Investment management and AI recommendations  |
| User Service       | http://localhost:8083  | User management and authentication            |
| AI Service         | http://localhost:8001  | AI-powered document verification and analysis |
| Prometheus         | http://localhost:9090  | Metrics collection                            |
| Grafana            | http://localhost:3000  | Metrics visualization (admin/admin)           |
| Jaeger             | http://localhost:16686 | Distributed tracing                           |

### Swagger UI

All services have interactive Swagger UI documentation available:

- KYC Service: http://localhost:8081/swagger-ui.html
- Investment Service: http://localhost:8082/swagger-ui.html
- User Service: http://localhost:8083/swagger-ui.html
- AI Service: http://localhost:8001/docs

### Test Token

For local development and testing, you can use the following JWT token which has access to all endpoints:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxOTE2MjM5MDIyLCJyb2xlcyI6WyJhZG1pbiIsInVzZXIiXX0.Ks0I-dCdjWUxJEwuGP0qlyYJGXXjUYlLCRwPIZXI5Ss
```

To use this token, add it to your API requests with the Authorization header:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlRlc3QgVXNlciIsImlhdCI6MTUxNjIzOTAyMiwiZXhwIjoxOTE2MjM5MDIyLCJyb2xlcyI6WyJhZG1pbiIsInVzZXIiXX0.Ks0I-dCdjWUxJEwuGP0qlyYJGXXjUYlLCRwPIZXI5Ss
```

## 📚 API Documentation

All services have comprehensive Swagger documentation. You can access the interactive Swagger UI for each service at the URLs listed above.

### Key Endpoints

#### KYC Service

- `GET /health` - Health check endpoint
- `POST /api/v1/auth/login` - Login endpoint
- `POST /api/v1/verifications` - Create a verification
- `GET /api/v1/verifications` - List all verifications
- `GET /api/v1/verifications/{id}` - Get verification by ID
- `PUT /api/v1/verifications/{id}` - Update verification
- `DELETE /api/v1/verifications/{id}` - Delete verification

#### Investment Service

- `POST /api/v1/investments` - Create a new investment
- `GET /api/v1/investments/{id}` - Get investment by ID
- `GET /api/v1/investments` - List all investments
- `PUT /api/v1/investments/{id}` - Update investment
- `DELETE /api/v1/investments/{id}` - Delete investment
- `POST /api/v1/portfolios` - Create a new portfolio
- `GET /api/v1/portfolios/{id}` - Get portfolio by ID
- `PUT /api/v1/portfolios/{id}` - Update portfolio
- `DELETE /api/v1/portfolios/{id}` - Delete portfolio
- `POST /api/v1/transactions` - Create a new transaction

#### User Service

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `GET /api/v1/users` - List all users
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register

#### AI Service

- `POST /api/v1/document/analyze` - Analyze and extract information from documents
- `POST /api/v1/image/classify` - Classify images using computer vision
- `POST /api/v1/text/analyze` - Analyze text for sentiment and entities

## ☁️ Cloud Deployment

SparkFund is designed to be deployed on major cloud providers with a focus on scalability, reliability, and security.

### Supported Cloud Providers

- **Amazon Web Services (AWS)**
  - EKS for Kubernetes orchestration
  - RDS for PostgreSQL databases
  - ElastiCache for Redis
  - S3 for document storage
  - CloudFront for content delivery
  - IAM for access management

- **Microsoft Azure**
  - AKS for Kubernetes orchestration
  - Azure Database for PostgreSQL
  - Azure Cache for Redis
  - Blob Storage for document storage
  - Azure CDN for content delivery
  - Azure AD for access management

- **Google Cloud Platform (GCP)**
  - GKE for Kubernetes orchestration
  - Cloud SQL for PostgreSQL databases
  - Memorystore for Redis
  - Cloud Storage for document storage
  - Cloud CDN for content delivery
  - IAM for access management

### Multi-Cloud Strategy

SparkFund supports a multi-cloud deployment strategy with:

- **Abstraction Layers**: Common interfaces for cloud-specific services
- **Portable Configurations**: Cloud-agnostic Kubernetes manifests
- **Consistent CI/CD**: Unified deployment pipelines across cloud providers
- **Cross-Cloud Networking**: Secure communication between services in different clouds

## 🏗 Infrastructure as Code

SparkFund uses Terraform for infrastructure provisioning and management.

### Terraform Modules

The `deploy/terraform` directory contains modular Terraform configurations:

- **Network**: VPC, subnets, security groups, and network policies
- **Kubernetes**: Managed Kubernetes clusters (EKS, AKS, GKE)
- **Databases**: Managed PostgreSQL instances
- **Caching**: Redis clusters
- **Storage**: Object storage buckets
- **IAM**: Identity and access management
- **Monitoring**: Prometheus, Grafana, and Jaeger infrastructure

### Environment Management

Terraform configurations are organized by environment:

- **Development**: For development and testing
- **Staging**: For pre-production validation
- **Production**: For production deployment

### Usage

```bash
# Initialize Terraform
cd deploy/terraform/environments/dev
terraform init

# Plan the deployment
terraform plan -out=tfplan

# Apply the changes
terraform apply tfplan
```

## 🐳 Containerization & Orchestration

SparkFund uses Docker for containerization and Kubernetes for orchestration.

### Docker

All services are containerized using Docker:

- **Base Images**: Alpine-based images for minimal footprint
- **Multi-stage Builds**: Optimized for small image sizes
- **Security Scanning**: Trivy and Clair integration for vulnerability scanning
- **Image Registry**: Support for Docker Hub, ECR, ACR, and GCR

### Kubernetes

Kubernetes is used for container orchestration:

- **Deployment Strategies**: Rolling updates, blue/green, and canary deployments
- **Auto-scaling**: Horizontal Pod Autoscaler (HPA) for dynamic scaling
- **Resource Management**: Resource requests and limits for all containers
- **Network Policies**: Secure communication between services
- **Service Mesh**: Istio integration for advanced traffic management

### Helm Charts

Helm is used for packaging and deploying Kubernetes applications:

- **Chart Structure**: Modular charts for each service
- **Value Overrides**: Environment-specific configurations
- **Dependencies**: Managed chart dependencies
- **Hooks**: Pre/post-install hooks for database migrations

### Usage

```bash
# Deploy using Helm
cd deploy/helm
helm install sparkfund ./sparkfund -f values-dev.yaml

# Upgrade an existing deployment
helm upgrade sparkfund ./sparkfund -f values-dev.yaml
```

## 💰 Cost Management

SparkFund integrates Kubecost for Kubernetes cost monitoring and optimization.

### Kubecost Features

- **Cost Allocation**: Track costs by namespace, deployment, and label
- **Cost Optimization**: Recommendations for right-sizing resources
- **Budget Alerts**: Notifications for cost anomalies
- **Chargeback Reports**: Cost allocation reports for internal billing
- **Savings Insights**: Identify opportunities for cost reduction

### Cost Optimization Strategies

- **Spot Instances**: Use spot/preemptible instances for non-critical workloads
- **Autoscaling**: Scale down during low-traffic periods
- **Resource Right-sizing**: Optimize CPU and memory requests
- **Storage Tiering**: Use appropriate storage classes for different data
- **Cluster Bin Packing**: Efficient node utilization

### Access Kubecost

When deployed, Kubecost is available at:

```
http://kubecost.your-domain.com
```

## 🧪 Testing

SparkFund includes comprehensive tests for all services:

### Running Tests

```bash
# Run all tests
go test ./tests/...

# Run tests for a specific service
go test ./services/kyc-service/...
go test ./services/investment-service/...
go test ./services/user-service/...
go test ./services/api-gateway/...
```

### Test Structure

Tests are organized by service and type:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test interactions between components
- **End-to-End Tests**: Test complete workflows

## 📊 Monitoring and Observability

SparkFund includes comprehensive monitoring and observability:

### Prometheus

Prometheus collects metrics from all services. Access the Prometheus UI at http://localhost:9090.

### Grafana

Grafana provides visualization of metrics with pre-configured dashboards. Access Grafana at http://localhost:3000 (username: admin, password: admin).

### Jaeger

Jaeger provides distributed tracing for request flows. Access the Jaeger UI at http://localhost:16686.

### Logging

All services use structured logging with consistent formats. Logs can be viewed using:

```bash
# View logs for all services
docker-compose logs

# View logs for a specific service
docker-compose logs [service-name]
```

## 🔧 Troubleshooting

If you encounter issues:

1. Check if all containers are running: `docker-compose ps`
2. View logs for a specific service: `docker-compose logs [service-name]`
3. Restart a specific service: `docker-compose restart [service-name]`
4. Restart all services: `docker-compose down && docker-compose up -d`

### Common Issues

- **Service Unavailable**: Check if the service container is running and healthy
- **Database Connection Issues**: Check if PostgreSQL is running and accessible
- **Authentication Failures**: Verify the JWT token is valid and not expired
- **API Gateway Routing Issues**: Check the Nginx configuration in the API Gateway service
- **Kubernetes Deployment Issues**: Check pod status and events with `kubectl describe pod`
- **Terraform Errors**: Check the Terraform state and logs for detailed error messages

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

<p align="center">
  Made with ❤️ by the SparkFund Team
</p>
