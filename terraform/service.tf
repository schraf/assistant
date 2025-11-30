
resource "google_cloud_run_v2_service" "assistant" {
  name     = "assistant"
  location = var.region

  ingress = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.assistant.email

    containers {
      name = "assistant"

      # Use latest tagged image from Artifact Registry, or custom image if specified
      image = var.initial_image != "" ? var.initial_image : "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.assistant.repository_id}/assistant:latest"

      ports {
        container_port = 8080
      }

      env {
        name  = "GOOGLE_CLOUD_PROJECT"
        value = var.project_id
      }

      env {
        name  = "API_TOKEN"
        value = var.api_token
      }

      env {
        name  = "CLOUD_RUN_JOB_NAME"
        value = google_cloud_run_v2_job.assistant.name
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
    google_artifact_registry_repository.assistant,
  ]
}

# Allow public access to the assistant cloud run service
resource "google_cloud_run_v2_service_iam_binding" "assistant_public_access" {
  location = google_cloud_run_v2_service.assistant.location
  project  = google_cloud_run_v2_service.assistant.project
  name     = google_cloud_run_v2_service.assistant.name
  role     = "roles/run.invoker"
  members = [
    "allUsers",
  ]
}
