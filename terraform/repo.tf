resource "google_sourcerepo_repository" "toto-config-api" {
  name = var.repository_name

  depends_on = [
    google_project_service.source_repo,
  ]
}

resource "google_cloudbuild_trigger" "trigger" {
  trigger_template {
    branch_name = "main"
    repo_name   = google_sourcerepo_repository.toto-config-api.name
  }

  filename = "cloudbuild.yaml"

  depends_on = [google_sourcerepo_repository.toto-config-api]
}
