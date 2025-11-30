
resource "google_artifact_registry_repository" "assistant" {
  location      = var.region
  repository_id = "assistant"
  description   = "Docker repository for Assistant container images"
  format        = "DOCKER"

  depends_on = [google_project_service.required_apis]
}

