
# Grant Cloud Build service account permission to push images to Artifact Registry
data "google_project" "project" {
  project_id = var.project_id
}

resource "google_artifact_registry_repository_iam_member" "cloudbuild_writer" {
  location   = var.region
  repository = google_artifact_registry_repository.assistant.repository_id
  role       = "roles/artifactregistry.writer"
  member     = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}

# Grant Cloud Build service account permission to deploy to Cloud Run
# Note: This uses project-level IAM to avoid circular dependencies
resource "google_project_iam_member" "cloudbuild_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = "serviceAccount:${data.google_project.project.number}@cloudbuild.gserviceaccount.com"
}

