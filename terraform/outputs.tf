output "service_account_email" {
  description = "Email of the created service account"
  value       = google_service_account.assistant.email
}

output "service_account_id" {
  description = "ID of the created service account"
  value       = google_service_account.assistant.id
}

output "service_account_name" {
  description = "Name of the created service account"
  value       = google_service_account.assistant.name
}

output "cloud_run_service_url" {
  description = "URL of the deployed Cloud Run service"
  value       = google_cloud_run_v2_service.assistant.uri
}

output "artifact_registry_repository" {
  description = "Full path to the Artifact Registry repository"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.assistant_repo.repository_id}"
}

output "container_image" {
  description = "Full container image path"
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.assistant_repo.repository_id}/assistant:latest"
}

output "cloud_run_job_name" {
  description = "Name of the Cloud Run Job for processing requests"
  value       = google_cloud_run_v2_job.assistant_job.name
}

output "api_gateway_url" {
  description = "URL of the API Gateway endpoint"
  value       = google_api_gateway_gateway.assistant_gateway.default_hostname
}

output "api_gateway_id" {
  description = "ID of the API Gateway"
  value       = google_api_gateway_gateway.assistant_gateway.gateway_id
}

output "api_key" {
  description = "API key for accessing the API Gateway (use this in the x-api-key header)"
  value       = google_apikeys_key.api_gateway_key.key_string
  sensitive   = true
}

output "api_key_name" {
  description = "Name of the API key resource"
  value       = google_apikeys_key.api_gateway_key.name
}

