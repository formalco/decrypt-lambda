output "decryptor_uri" {
  description = "The /decrypt endpoint to set as the Formal encryption key's decryptor URI."
  value       = "${google_cloud_run_v2_service.decryptor.uri}/decrypt"
}

output "service_account_email" {
  description = "Service account the decryptor runs as."
  value       = google_service_account.decryptor.email
}
