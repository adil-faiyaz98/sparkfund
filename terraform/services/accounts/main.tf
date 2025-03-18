resource "kubernetes_deployment" "accounts_service" {
  metadata {
    name = "accounts-service"
    namespace = var.namespace
    labels = {
      app = "accounts-service"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        app = "accounts-service"
      }
    }

    template {
      metadata {
        labels = {
          app = "accounts-service"
        }
      }

      spec {
        container {
          name  = "accounts-service"
          image = "${var.ecr_repository_url}:${var.image_tag}"
          port {
            container_port = 8080
            name          = "http"
          }

          env {
            name  = "DB_HOST"
            value = var.db_host
          }

          env {
            name  = "DB_PORT"
            value = var.db_port
          }

          env {
            name  = "DB_USER"
            value = var.db_user
          }

          env {
            name  = "DB_PASSWORD"
            value_from {
              secret_key_ref {
                name = "accounts-db-secret"
                key  = "password"
              }
            }
          }

          env {
            name  = "DB_NAME"
            value = var.db_name
          }

          resources {
            limits = {
              cpu    = "500m"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "256Mi"
            }
          }

          readiness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 5
            period_seconds       = 10
          }

          liveness_probe {
            http_get {
              path = "/health"
              port = 8080
            }
            initial_delay_seconds = 15
            period_seconds       = 20
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "accounts_service" {
  metadata {
    name      = "accounts-service"
    namespace = var.namespace
  }

  spec {
    selector = {
      app = "accounts-service"
    }

    port {
      port        = 80
      target_port = "http"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_secret" "accounts_db_secret" {
  metadata {
    name      = "accounts-db-secret"
    namespace = var.namespace
  }

  data = {
    password = var.db_password
  }
} 