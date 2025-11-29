terraform {
  required_version = ">= 1.0"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}

provider "google-beta" {
  project             = var.project_id
  region              = var.region
  user_project_override = true
}

# Enable required APIs
resource "google_project_service" "required_apis" {
  for_each = toset([
    "cloudresourcemanager.googleapis.com",
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "apigateway.googleapis.com",
    "servicemanagement.googleapis.com",
    "servicecontrol.googleapis.com",
    "apikeys.googleapis.com",
  ])

  project = var.project_id
  service = each.value

  disable_dependent_services = true
  disable_on_destroy = true
}

# Create Artifact Registry repository for container images
resource "google_artifact_registry_repository" "assistant_repo" {
  location      = var.region
  repository_id = "assistant"
  description   = "Docker repository for Assistant container images"
  format        = "DOCKER"

  depends_on = [google_project_service.required_apis]
}

# Cloud Run service
resource "google_cloud_run_v2_service" "assistant" {
  name     = "assistant"
  location = var.region

  # API Gateway requires the service to be publicly accessible
  # Access is restricted via IAM to only allow the API Gateway service account
  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.assistant.email

    containers {
      name = "assistant"
      # Use latest tagged image from Artifact Registry, or custom image if specified
      image = var.initial_image != "" ? var.initial_image : "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.assistant_repo.repository_id}/assistant:latest"

      ports {
        container_port = 8080
      }

      env {
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project_id
      }

      env {
        name  = "CLOUD_RUN_JOB_NAME"
        value = google_cloud_run_v2_job.assistant_job.name
      }

      env {
        name  = "CLOUD_RUN_JOB_REGION"
        value = var.region
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "1Gi"
        }
      }
    }

    timeout = "60s"
  }

  depends_on = [
    google_project_service.required_apis,
    google_artifact_registry_repository.assistant_repo,
  ]
}

# Create Cloud Run Job for processing requests
resource "google_cloud_run_v2_job" "assistant_job" {
  name     = "assistant-job"
  location = var.region

  template {
    template {
      service_account = google_service_account.assistant.email

      containers {
        image = var.initial_image != "" ? var.initial_image : "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.assistant_repo.repository_id}/assistant:latest"

        command = ["/job"]

        env {
          name  = "GOOGLE_API_KEY"
          value = var.google_api_key
        }

        env {
          name  = "TELEGRAPH_API_KEY"
          value = var.telegraph_api_key
        }

        env {
          name  = "TELEGRAPH_AUTHOR_NAME"
          value = var.telegraph_author_name
        }

        env {
          name  = "MAIL_SMTP_SERVER"
          value = var.smtp_hostname
        }

        env {
          name  = "MAIL_SMTP_PORT"
          value = var.smtp_port
        }

        env {
          name  = "MAIL_SENDER_EMAIL"
          value = var.mail_sender_email
        }

        env {
          name  = "MAIL_SENDER_PASSWORD"
          value = var.mail_sender_password
        }

        env {
          name  = "MAIL_RECIPIENT_EMAIL"
          value = var.mail_recipient_email
        }

        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = var.project_id
        }

        resources {
          limits = {
            cpu    = "1"
            memory = "1Gi"
          }
        }
      }

      timeout     = "3600s"
      max_retries = 2
    }
  }

  depends_on = [
    google_project_service.required_apis,
    google_artifact_registry_repository.assistant_repo,
    google_service_account.assistant,
  ]
}

# API Gateway API configuration
resource "google_api_gateway_api" "assistant_api" {
  provider     = google-beta
  api_id       = "assistant-api"
  display_name = "Assistant API"
  project      = var.project_id

  depends_on = [google_project_service.required_apis]
}

# API Gateway API Config - uses the OpenAPI spec
resource "google_api_gateway_api_config" "assistant_api_config" {
  provider    = google-beta
  api         = google_api_gateway_api.assistant_api.api_id
  api_config_id = "assistant-api-config-${formatdate("YYYYMMDDhhmmss", timestamp())}"

  openapi_documents {
    document {
      path     = "openapi.yaml"
      contents = base64encode(templatefile("${path.module}/../openapi.yaml", {
        CLOUD_RUN_SERVICE_URL = google_cloud_run_v2_service.assistant.uri
        API_ID                = "assistant-api"
      }))
    }
  }

  gateway_config {
    backend_config {
      google_service_account = google_service_account.api_gateway.email
    }
  }

  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    google_api_gateway_api.assistant_api,
    google_cloud_run_v2_service.assistant,
    google_service_account.api_gateway,
  ]
}

# API Gateway Gateway
resource "google_api_gateway_gateway" "assistant_gateway" {
  provider   = google-beta
  region     = var.region
  project    = var.project_id
  api_config = google_api_gateway_api_config.assistant_api_config.id
  gateway_id = "assistant-gateway"

  depends_on = [google_api_gateway_api_config.assistant_api_config]
}

# Restrict Cloud Run service access to only the API Gateway service account
resource "google_cloud_run_v2_service_iam_binding" "api_gateway_only" {
  location = google_cloud_run_v2_service.assistant.location
  name     = google_cloud_run_v2_service.assistant.name
  role     = "roles/run.invoker"

  members = [
    "serviceAccount:${google_service_account.api_gateway.email}",
  ]

  depends_on = [
    google_cloud_run_v2_service.assistant,
    google_service_account.api_gateway,
  ]
}

# Create API key for API Gateway
resource "google_apikeys_key" "api_gateway_key" {
  provider = google-beta
  name     = "assistant-api-key"
  project  = var.project_id

  display_name = "Assistant API Gateway Key"

  restrictions {
    # Restrict the API key to API Gateway only
    api_targets {
      service = "apigateway.googleapis.com"
    }
  }

  depends_on = [
    google_project_service.required_apis,
    google_api_gateway_api.assistant_api,
  ]
}

