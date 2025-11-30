# Create the service account
resource "google_service_account" "assistant" {
  account_id   = "assistant"
  display_name = "Assistant Service Account"
  description  = "Service account for Assistant Cloud Run service"
}

# Grant the service account permission to use Gemini API
resource "google_project_iam_member" "assistant_aiplatform" {
  project = var.project_id
  role    = "roles/aiplatform.user"
  member  = "serviceAccount:${google_service_account.assistant.email}"
}

# Grant the service account permission to execute Cloud Run Jobs
# Note: run.jobsExecutor allows executing jobs, but we also need run.admin or run.developer
# to use RunJob API with overrides
resource "google_project_iam_member" "assistant_run_jobs_executor" {
  project = var.project_id
  role    = "roles/run.jobsExecutor"
  member  = "serviceAccount:${google_service_account.assistant.email}"
}

# Grant the service account permission to run jobs with overrides
# This is needed for the RunJob API call with environment variable overrides
resource "google_project_iam_member" "assistant_run_developer" {
  project = var.project_id
  role    = "roles/run.developer"
  member  = "serviceAccount:${google_service_account.assistant.email}"
}

