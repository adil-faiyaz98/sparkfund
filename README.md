# SparkFund Platform [ POC - AI Integration ]

SparkFund is a **Proof of Concept** for AI-powered investment recommendations, KYC verification, and user management capabilities. Built with a microservices architecture, it provides a foundation for exploring how AI could enhance financial applications.

## Table of Contents

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
- [Images](#images)
- [Troubleshooting](#-troubleshooting)

## Features

- **Microservices Architecture**: Scalable, maintainable, and resilient system design
- **AI-Powered Investment Recommendations**: Proof-of-concept algorithms for personalized investment advice
- **KYC Verification**: Simulated customer onboarding with mock AI-powered document verification
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

1. **API Gateway**: Routes requests to appropriate services, handles authentication, rate limiting, load balancing, request validation, security, A/B testing, canary deployments and more
2. **KYC Service**: Handles customer verification and onboarding with AI-powered document verification
3. **Investment Service**: Manages investments and provides AI-powered recommendations
4. **User Service**: Manages user accounts, authentication, and profiles
5. **AI Service**: Provides AI capabilities for document verification, facial recognition, and investment analysis
6. **Supporting Infrastructure**: PostgreSQL, Redis, Prometheus, Grafana, and Jaeger

## Services

| Service            | Description                                                                                                                                                                                                     |
| ------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| API Gateway        | Routes requests to appropriate services, handles authentication, rate limiting, and request validation                                                                                                          |
| KYC Service        | Handles customer verification and onboarding with API-based integration to AI services [Uses API forwarding to an external AI service; most endpoints return simulated responses]                               |
| Investment Service | Manages investments, portfolios, and provides AI-powered recommendations and market analysis [Contains well-defined interfaces for AI models but primarily uses simulated data rather than real trained models] |
| User Service       | Manages user accounts, authentication, profiles, and permissions                                                                                                                                                |
| AI Service         | Provides AI capabilities for document verification, facial recognition, NLP, and investment analysis [POC implementation with placeholder code for model integration]                                           |

## AI Capabilities

SparkFund integrates AI capabilities across its services as a **Proof of Concept (POC)**. The core microservices functionality is implemented, but the AI features are primarily **simulations** demonstrating how machine learning could enhance financial services.

### KYC Service AI (POC Implementation)

- **Document Verification**: Analysis of identity documents for authenticity and data extraction

  - _Integration Pattern_: API forwarding to an external AI service
  - _Implementation Status_: **API Forwarding with Mock Fallbacks** - Most endpoints call an external AI service, with local mock responses as fallbacks
  - _Production Readiness_: **Low** - Requires actual AI service with trained models for document verification

- **Facial Recognition**: Face matching between selfie and ID document

  - _Integration Pattern_: External API call to dedicated AI service
  - _Implementation Status_: **API Integration** - Requests forwarded to external service endpoint `/api/v1/face/match-base64`
  - _Production Readiness_: **Low** - Dependent on external service with actual face recognition models

- **Risk Analysis**: Assessment of customer risk profile

  - _Integration Pattern_: External API call to dedicated AI service
  - _Implementation Status_: **API Integration** - Requests forwarded to external service endpoint `/api/v1/risk/analyze`
  - _Production Readiness_: **Low** - Dependent on external service with actual risk assessment models

- **Anomaly Detection**: Identification of suspicious patterns in verification requests
  - _Integration Pattern_: External API call to dedicated AI service
  - _Implementation Status_: **API Integration** - Requests forwarded to external service endpoint `/api/v1/anomaly/detect`
  - _Production Readiness_: **Low** - Dependent on external service with actual anomaly detection models

### Investment Service AI (POC Implementation)

- **Recommendation System**: Personalized investment recommendations based on user profile and market conditions

  - _Models_: Interfaces defined for collaborative filtering, content-based filtering, and hybrid recommendation systems
  - _Implementation Status_: **Simulation** - API endpoints return realistic sample responses that mimic model predictions
  - _Production Readiness_: **Low** - Would require real model training, evaluation, and deployment

- **Natural Language Processing**: Analysis of market news for investment signals

  - _Models_: References to BERT-based sentiment analysis, named entity recognition, and topic modeling
  - _Implementation Status_: **Placeholder** - Text processing pipeline defined but not connected to actual models
  - _For use in Actual projects\*\*: - integrate with actual NLP models for more accurate analysis_

- **Time Series Forecasting**: Price prediction and market movement analysis

  - _Models_: Interfaces for ARIMA, LSTM, and Prophet models
  - _Implementation Status_: **Simulation** - Returns sample forecasts based on historical patterns
  - _For use in Actual projects\*\*: - integrate with actual time series forecasting models for more accurate predictions_

- **Reinforcement Learning**: Portfolio optimization and investment strategy learning

  - _Models_: Scaffold code for Deep Q-Networks (DQN) and Proximal Policy Optimization (PPO)
  - _Implementation Status_: **Conceptual** - Framework defined but lacks actual RL implementation
  - _For use in Actual projects\*\*: - integrate with actual RL algorithms for more sophisticated investment strategies_

- **Anomaly Detection**: Fraud detection and security monitoring
  - _Models_: References to Isolation Forest, One-Class SVM, and Autoencoder-based anomaly detection
  - _Implementation Status_: **Simulation** - Provides mock anomaly scores without real detection
  - _For use in Actual projects\*\*: - integrate with actual anomaly detection models for enhanced security_

### AI Service (Central AI Capabilities)

- **Document Analysis**: Extract information from documents using computer vision and NLP

  - _Models_: CRAFT for text detection, BERT for information extraction
  - _Implementation_: End-to-end document processing pipeline
  - _Endpoints_: `/api/v1/document/analyze`

- **Image Classification**: Classify images for various purposes

  - _Models_: EfficientNet, ResNet for image classification
  - _Implementation_: Image processing and classification endpoints
  - _Endpoints_: `/api/v1/image/classify`

- **Text Analysis**: Sentiment analysis, entity extraction, and text classification
  - _Models_: DistilBERT for sentiment analysis, SpaCy for NER
  - _Implementation_: Text processing endpoints with multiple analysis options
  - _Endpoints_: `/api/v1/text/analyze`

### POC Implementation Notes

The AI capabilities are implemented as a proof of concept with the following characteristics:

1. **API Endpoints**: All AI endpoints are fully defined and documented in Swagger
2. **Integration Patterns**:
   - KYC Service uses API forwarding to an external AI service
   - Investment Service uses direct simulation within the service
3. **Sample Responses**: Many endpoints return realistic sample responses based on the expected model output
4. **Service Configuration**: Comprehensive configuration for AI service integration is defined but points to placeholder services
5. **Demonstration Ready**: The system demonstrates the potential of AI in financial services
6. **Production Path**: Would require significant additional work to implement or integrate with real AI models

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### Clone the Repository

```bash
git clone https://github.com/adil-faiyaz98/sparkfund.git
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

## API Documentation

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

## Cloud Deployment

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

## Infrastructure as Code

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

## Containerization & Orchestration

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

## Cost Management

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

## Testing

SparkFund includes comprehensive tests for all services:

### Running Tests

- Not going to run tests for now. But, some really important security related tests are in api-gateway. As for the others, I would recommend the builders for in-depth testing.

## Monitoring and Observability

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

## Images

### SparkFund Logo

![SparkFund Logo](assets/images/sparkfund-logo.svg)

### Architecture Diagram

![SparkFund Architecture](assets/images/2025-04-07 17_56_59.png)

### Swagger UI

![Swagger UI](assets/images/2025-04-07 17_56_59-Greenshot.png)

### Monitoring

#### Prometheus

![Prometheus Dashboard](assets/images/2025-04-07 17_30_25.png)

#### Grafana

![Grafana Dashboard](assets/images/2025-04-07 17_56_54.png)

#### Jaeger

![Jaeger UI](assets/images/2025-04-07 17_57_34.png)

## Troubleshooting

If you encounter issues:

## ChatGPT for Help. Otherwise:

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

---
