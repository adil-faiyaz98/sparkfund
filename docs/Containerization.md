# Contarinerization Technologies

## Docker

- Docker helps to containerize the application so that it is always run in the same environment.
- It is a platform that allows you to develop, ship, and run applications inside containers.
- Containers are lightweight, standalone, and executable software packages that include everything needed to run a piece of software, including the code, runtime, system tools, libraries, and settings.

## Kubernetes

- Kubernetes is an open-source system for automating deployment, scaling, and management of containerized applications. It is designed to work with Docker containers.
- It is a container orchestration platform that automates the deployment, scaling, and management of containerized applications.
- It provides a framework for running distributed systems resiliently.
- It is designed to automate the deployment, scaling, and management of containerized applications.
- It is a platform that allows you to run, manage, and scale containerized applications.
- It is a system for automating the deployment, scaling, and management of containerized applications.
- It is a platform that allows you to run, manage, and scale containerized applications.

## What is Terraform ?

- Terraform is an open-source infrastructure as code (IaC) software tool created by HashiCorp. Users define and provision data center infrastructure using a declarative configuration language known as HashiCorp Configuration Language (HCL), or optionally JSON.
- Terraform generates an execution plan describing what it will do to reach the desired state, and then executes it to build the described infrastructure. The infrastructure Terraform can manage includes low-level components such as compute instances, storage, and networking, as well as high-level components such as DNS entries and SaaS features.
- Terraform can manage both existing service providers and custom in-house solutions.
- Configuration files describe to Terraform the components needed to run a single application or your entire datacenter. Terraform generates an execution plan describing what it will do to reach the desired state, and then executes it to build the described infrastructure.

## Terraform and Kubernetes

- Terraform is used to provision infrastructure, while Kubernetes is used to manage containers. Terraform can be used to provision a Kubernetes cluster, and then Kubernetes can be used to manage the containers that run on that cluster.

## Kubernetes and Docker

- Docker containers are used to package and run applications. Kubernetes is used to manage and orchestrate Docker
- Orchestration basically is achieved by Kubernetes by managing the deployment, scaling, and management of containerized applications.
- Kubernetes is a platform that allows you to run, manage, and scale containerized applications.
- It uses Docker containers to package and run applications.
- Kubernetes is a container orchestration platform that automates the deployment, scaling, and management of containerized applications.
- It is designed to work with Docker containers.
- Kubernetes is a platform that allows you to run, manage, and scale containerized applications.

## Kubernetes and Ansible

- Kubernetes is used to manage and orchestrate containers, while Ansible is used to automate the provisioning and configuration of infrastructure.
- The two tools can be used together to automate the deployment and management of containerized applications.
- Kubernetes can be used to manage the containers that are provisioned by Ansible.
- Ansible can be used to provision the infrastructure that Kubernetes runs on.

## Ansible and Terraform

-

## Kubernetes Clusters (Types)

- Kubernetes clusters are typically of two types: on-premises and cloud-based.
- On-premises clusters are deployed and managed within an organization's own data center, while cloud-based clusters are deployed and managed by a cloud service provider.
- On-premises clusters can be further divided into two types: bare-metal and virtualized.
- Bare-metal clusters are deployed on physical servers, while virtualized clusters are deployed on virtual machines.
- Cloud-based clusters can be further divided into two types: managed and self-managed.
- Managed clusters are deployed and managed by a cloud service provider, while self-managed clusters are deployed and managed by an organization.

- ## What's the difference between AWS EKS vs EC2 instances ?
- AWS EKS (Elastic Kubernetes Service) is a managed Kubernetes service provided by AWS, while EC2 instances are virtual servers provided by AWS.
- The main difference between the two is that EKS handles the control plane of the Kubernetes cluster, while EC2 instances are used to run the worker nodes.
- EKS is a managed service, which means that AWS handles the provisioning, scaling, and management of the control plane, while EC2 instances are self-managed, which means that the organization is responsible for the provisioning, scaling, and management of the worker nodes.
- EKS is a fully managed service, which means that AWS handles all aspects of the Kubernetes cluster, including the control plane and the worker nodes.
- The computer power of EC2 instances is determined by the instance type, while the computer power of EKS is determined by the number of worker nodes and the instance type of the worker nodes.
