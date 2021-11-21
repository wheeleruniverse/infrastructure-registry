
output "service_account_json" {
  sensitive = true
  value     = google_service_account_key.this.private_key
}
