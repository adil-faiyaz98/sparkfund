# NGNIX

## What is NGINX and Why Use It?

- NGINX is a popular open-source web server and reverse proxy server that is known for its high performance, stability, and rich feature set. It is widely used for serving static content, load balancing, caching, and acting as a reverse proxy for other servers.
- It sits in between the clients and the microservices deployed in the cloud. It helps in load balancing, caching, and providing security features.
- It's main purpose is to act as a reverse proxy, load balancer, and caching server.

## How to Install NGINX?

- The installation process for NGINX can vary depending on the operating system you are using. Below are the general steps for installing NGINX on different platforms:

````bash
    - **Ubuntu/Debian:**
        ```sh
        sudo apt update
        sudo apt install nginx
        ```
    - **CentOS/RHEL:**
        ```sh
        sudo yum install epel-release
        sudo yum install nginx
        ```
    - **macOS (using Homebrew):**
        ```sh
        brew install nginx
        ```
    - **Windows:**
        - Download the NGINX installer from the [official website](https://nginx.org/en/download.html).
        - Follow the installation instructions provided in the installer.
````

## What configuration should we give it and where do we install it ? (In a microservices architecture)

In a microservices architecture, NGINX is typically used as a reverse proxy, load balancer, and caching server. Here are some common configurations and considerations for deploying NGINX in such an environment:

- **Reverse Proxy Configuration:**
- Configure NGINX to forward requests to the appropriate microservices based on the URL path or domain. For example:

```bash
server {
    listen 80;
    server_name example.com;
        location /api/ {
            proxy_pass http://api-service:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
                 }
        location / {
            proxy_pass http://frontend-service:8080;
                proxy_set_header Host $host;
                proxy_set_header X-Real-IP $remote_addr;
                proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $scheme;
        }
}
```

- **Load Balancing Configuration:** Configure NGINX to distribute incoming traffic across multiple instances of a microservice. For example:

```bash
upstream api-service {
    server api-service1:8080;
    server api-service2:8080;
    server api-service3:8080;
}
```

# Kubernetes Ingress Controller vs NGINX

## what is the difference between Kubernetes Ingress Controller and NGINX?

- The Kubernetes Ingress Controller and NGINX serve similar purposes in a Kubernetes environment, but they have distinct roles and functionalities. Here’s a breakdown of the differences:

- **Kubernetes Ingress Controller:**
  - **Purpose:** The Ingress Controller is a Kubernetes resource that manages external access to the services in a cluster, typically HTTP and HTTPS routes.
  - **Functionality:** It acts as a reverse proxy and load balancer, routing traffic to the appropriate services based on the rules defined in Ingress resources.
  - **Configuration:** Ingress resources are defined using YAML files and applied to the Kubernetes cluster. These resources specify the routing rules and backend services.
  - **Flexibility:** Ingress controllers can be customized and extended with various plugins and annotations to support additional features.
  - **Examples:** NGINX Ingress Controller, Traefik, HAProxy Ingress Controller.
  - **Deployment:** Typically deployed as a Kubernetes Deployment or DaemonSet, and managed by Kubernetes.
  - **Configuration Management:** Managed through Kubernetes manifests and annotations.
  - **Scalability:** Scales automatically with the number of Ingress resources and services.
