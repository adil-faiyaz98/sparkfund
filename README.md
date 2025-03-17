# AegisFinance

Kubernetes-Based Microservices with AI/ML & Big Data

## Overview

This project sets up a Kubernetes-based microservices architecture with AI/ML pipelines, CI/CD automation, and Big Data processing, all within the free tier indefinitely or at least for one year across cloud providers.

## Tech Stack

### Cloud Services

- **AWS Elastic Kubernetes Service**: Managed Kubernetes cluster
- **Google Cloud Run**: Serverless API hosting (optional)
- **Containerization**: Docker
- **Terraform**: Infrastructure as Code (IaC) & Kubernetes Manifests to deploy microservices to AWS EKS

### Microservices Architecture

- **Language**: Go (GIN framework)
- **Total Services**: Microservices
- **Endpoints per Service**: ~5 endpoints each
- **Database** : PostgreSQL
- **Automated Testing** - Automated with Ginkgo and Gomega
- **CI/CD**: GitHub Runners (self-hosted)

### Data & AI/ML

- **Big Data Processing**: Google BigQuery for Analytics, Data Warehousing, and Machine Learning
- **AI Models**: Fraud detection, anomaly detection, credit card recommendations
- **AI Model Training**: Training on Jupyter Notebook, having GPU/CPU instances on AWS SageMaker
- **AI Model Deployment**: Hosted on AWS SageMaker Endpoint for real-time inference, AWS Lambda, or Vertex AI
- **ETL Processing and Data Pipelines**: AWS Glue for automated data processing or Kubernetes CronJobs

### Monitoring & Logging

- **Prometheus**: Monitoring and alerting system
- **Grafana**: Visualization and monitoring dashboard
- **ELK Stack**: Elasticsearch, Logstash, and Kibana for logging and monitoring
- **Jaeger**: Distributed tracing system

### Load Balancing & Security

- **[Option1]NGINX** - Will be used as a reverse proxy and load balancer as a standalone service in Kubernetes. - It will also be used to manage SSL/TLS certificates for secure communication. - NGINX will be configured to handle API rate limiting, caching, and request throttling and request logging - It will also be used for security ( JWT validation ) - It will be used to route requests to the appropriate microservice based on the URL path.

- **[Option2]Kubernetes Ingress Controller** - Ingress manages external requests to the services in the cluster dynamically. It allows you to configure routing rules, SSL termination, and load balancing. - It is more native to Kubernetes and integrates with cert-manager for HTTPS. - It is more scalable and powerful than NGINX and integrates better with cloud-native load balancing.

## Setup Instructions

### 1. Install Dependencies

```bash
# Install AWS CLI
curl "https://awscli.amazonaws.com/AWSCLIV2.pkg" -o "AWSCLIV2.pkg"
sudo installer -pkg AWSCLIV2.pkg -target /

# Install Terraform
wget https://releases.hashicorp.com/terraform/X.X.X/terraform_X.X.X_linux_amd64.zip
unzip terraform_X.X.X_linux_amd64.zip
mv terraform /usr/local/bin/

# Install kubectl
aws eks update-kubeconfig --name my-cluster

# Install Docker
sudo apt-get update && sudo apt-get install docker.io -y

# Install Minikube for local Kubernetes testing
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin/
```

### 2. Setup AWS EKS Cluster

```bash
# Create an EKS cluster
aws eks create-cluster --name free-tier-cluster --role-arn arn:aws:iam::<YOUR_ACCOUNT_ID>:role/EKSClusterRole --resources-vpc-config subnetIds=<SUBNET_ID>,securityGroupIds=<SG_ID>
```

### 3. Deploy Microservices to Kubernetes

```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

### 4. Deploy AI/ML Model to AWS SageMaker

```bash
# Train and deploy AI Model on SageMaker
aws sagemaker create-training-job --training-job-name fraud-model --algorithm-specification TrainingImage=<YOUR_ALGO_IMAGE> --role-arn arn:aws:iam::<YOUR_ACCOUNT_ID>:role/SageMakerRole
```

### 5. Deploy PostgreSQL Database in Kubernetes

```bash
kubectl apply -f k8s/postgres.yaml
```

### 6. Deploy CI/CD Pipelines using GitHub Actions

```yaml
name: Deploy to EKS

on:
  push:
    branches:
      - main

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v2
      - name: Setup Terraform
        uses: hashicorp/setup-terraform@v1
      - name: Apply Terraform
        run: terraform apply -auto-approve
```
