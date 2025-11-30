
resource "google_cloud_run_v2_job" "assistant" {
  name     = "assistant-job"
  location = var.region

  template {
    template {
      service_account = google_service_account.assistant.email

      containers {
        image = var.initial_image != "" ? var.initial_image : "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.assistant.repository_id}/assistant:latest"

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
    google_artifact_registry_repository.assistant,
    google_service_account.assistant,
  ]
}

