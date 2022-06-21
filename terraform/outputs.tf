
output "config-api_http_url" {
  value = module.cloud_run_config_api_http.url
}

output "repo_url" {
  value = google_sourcerepo_repository.toto-config-api.url
}