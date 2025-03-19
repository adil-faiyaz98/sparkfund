# Accounts Service

The Accounts Service is a microservice responsible for managing user accounts in the Money Pulse application. It provides functionality for creating, retrieving, updating, and deleting accounts, as well as managing account balances.

## Features

- Create new accounts
- Retrieve account details
- Update account information
- Delete accounts
- Get all accounts for a user
- Update account balances
- Get account by account number

## Architecture

The service follows a hexagonal architecture pattern with the following components:

- **Domain**: Core business logic and entities
- **Application**: Use cases and business rules
- **Ports**: Interfaces for external communication
- **Adapters**: Implementations of ports (gRPC, PostgreSQL)

## Prerequisites

- Go 1.21 or later
- PostgreSQL 14 or later
- Protocol Buffers compiler (protoc)
- Make

## Getting Started

1. Install dependencies:
   ```bash
   make deps
   ```

2. Generate protobuf code:
   ```bash
   make proto
   ```

3. Build the service:
   ```bash
   make build
   ```

4. Run the service:
   ```bash
   make run
   ```

## Configuration

The service can be configured using environment variables:

- `DB_DSN`: PostgreSQL connection string
- `PORT`: gRPC server port (default: 50051)

## API Documentation

The service exposes a gRPC API with the following methods:

- `CreateAccount`: Create a new account
- `GetAccount`: Get account by ID
- `GetUserAccounts`: Get all accounts for a user
- `UpdateAccount`: Update account details
- `DeleteAccount`: Delete an account
- `GetAccountByNumber`: Get account by account number
- `UpdateBalance`: Update account balance

## Development

### Running Tests

```bash
make test
```

### Linting

```bash
make lint
```

## Deployment

The service can be deployed using Kubernetes and Helm:

1. Build the Docker image:
   ```bash
   docker build -t accounts-service:latest .
   ```

2. Deploy to Kubernetes:
   ```bash
   kubectl apply -k k8s/overlays/development
   ```

## Infrastructure

The service infrastructure is managed using Terraform:

- EKS cluster for container orchestration
- RDS instance for PostgreSQL database
- Security groups and IAM roles
- VPC and subnet configuration

### Development Environment

```bash
cd terraform/environments/development
terraform init
terraform plan
terraform apply
```

### Production Environment

```bash
cd terraform/environments/production
terraform init
terraform plan
terraform apply
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 