# AWS EKS

## What is it?

- AWS EKS stands for Amazon Elastic Kubernetes Service.
- It is a managed service that makes it easy to run Kubernetes on AWS without needing to install and operate your own Kubernetes control plane or nodes.
- EKS is fully integrated with other AWS services, such as IAM, VPC, and CloudWatch, making it easier to manage and secure your Kubernetes clusters.

## Infrastructure

- EKS runs on AWS infrastructure and provides a highly available and scalable Kubernetes control plane.
- The control plane is managed by AWS, so you don't have to worry about updates, patches, or maintenance.
- You can run your worker nodes on AWS EC2 instances, Fargate, or any other Kubernetes-compatible infrastructure.
- EKS supports both public and private clusters, allowing you to choose the best networking model for your use case.
- EKS also supports multi-region clusters, allowing you to deploy your applications across multiple AWS regions for high availability and disaster recovery.
- EKS supports both Linux and Windows worker nodes, allowing you to run both Linux and Windows applications in the same cluster.
- EKS supports both Amazon ECR and Docker Hub for container images.

- EKS runs on CPUs mostly and are a level up from EC2 instances. It is used to host the microservices that were developed

## EKS Architecture

- EKS clusters consist of two main components: the control plane and the data plane.
- The control plane is responsible for managing the Kubernetes API server, etcd, and other Kubernetes components.
- The data plane consists of the worker nodes, which run the containerized applications.
- EKS provides a managed control plane, so you don't have to worry about managing the Kubernetes API server, etcd, and other Kubernetes components.
- You can run your worker nodes on AWS EC2 instances, Fargate, or any other Kubernetes-compatible infrastructure.
- EKS supports both public and private clusters, allowing you to choose the best networking model for your use case.
- EKS also supports multi-region clusters, allowing you to deploy your applications across multiple AWS regions for high availability and disaster recovery.

## What is Service Discovery?

- Service discovery is the process of automatically detecting and registering the services available in a network. It allows applications to find and communicate with each other without hardcoding IP addresses or DNS names.

## What is a Service Mesh

- A service mesh is a dedicated infrastructure layer for handling service-to-service communication. It provides features such as traffic management, security, observability, and policy enforcement.
- Popular service meshes include Istio, Linkerd, and Consul Connect.
- Service meshes can be deployed as a sidecar proxy alongside your application containers, intercepting and managing network traffic.
- They provide advanced features like traffic routing, load balancing, circuit breaking, and more.
- Service meshes can be used with any programming language or framework, making them language-agnostic.
- They can be deployed on any infrastructure, including on-premises, in the cloud, or in a hybrid environment.
- Service meshes can be used with any service discovery mechanism, including DNS, Consul, or Kubernetes.

# Comparison between EKS and SageMaker

- Both AWS EKS and Amazon SageMaker are services provided by AWS, but they serve different purposes and have different use cases. Here's a comparison between both the services:

## AWS EKS

### Purpose

- AWS EKS (Elastic Kubernetes Service) is a managed service that makes it easy to run Kubernetes on AWS without needing to install, operate, and maintain your own Kubernetes control plane. It enables you to run Kubernetes applications on AWS infrastructure.

### Use Cases

- Running containerized applications
- Microservices architecture
- Continuous integration and continuous deployment (CI/CD)
- Hybrid and multi-cloud deployments
- Serverless applications with AWS Lambda and Fargate

## SageMaker

### Purpose

- SageMaker is a fully managed service that provides every developer and data scientist with the ability to build, train, and deploy machine learning (ML) models quickly. It is designed to help you build, train, and deploy machine learning models at scale.

### Use Cases

-
