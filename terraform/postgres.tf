

provider "google-beta" {
  project = var.project
  region = var.region
  credentials = base64decode(google_service_account_key.postgres_key.private_key)
}

resource "google_service_account" "postgres" {
  account_id = "postgres"
  display_name = "Postgres Service Account"
  depends_on = [
    google_project_iam_member.owner,
    google_project_service.postgres,
  ]
}

resource "google_project_iam_member" "service_account_postgres_admin" {
  role   = "roles/editor"
  member = "serviceAccount:${google_service_account.postgres.email}"
  project = var.project
}

resource "google_service_account_key" "postgres_key" {
  service_account_id = google_service_account.postgres.name
}

resource "google_sql_database_instance" "postgresql" {
  name = "config-api-postgres"
  project = var.project
  region = var.region

  database_version = var.db_version
  deletion_protection = false

  settings {
    tier = var.db_tier
    activation_policy = var.db_activation_policy
    disk_autoresize = var.db_disk_autoresize
    disk_size = var.db_disk_size
    disk_type = var.db_disk_type
    pricing_plan = var.db_pricing_plan

    location_preference {
      zone = var.zone
    }

    maintenance_window {
      day  = "7"  # sunday
      hour = "3" # 3am
    }

    backup_configuration {
      enabled = true
      start_time = "00:00"
    }

    ip_configuration {
      ipv4_enabled = "true"
      authorized_networks {
        value = var.db_instance_access_cidr
      }
    }
  }
}
resource "google_sql_database" "postgresql_db" {
  name = var.db_name
  project = var.project
  instance = google_sql_database_instance.postgresql.name
}

# create user
resource "random_id" "user_password" {
  byte_length = 8
}

resource "google_sql_user" "postgresql_user" {
  name = var.db_user_name
  project  = var.project
  instance = google_sql_database_instance.postgresql.name
  password = var.db_user_password == "" ? random_id.user_password.hex : var.db_user_password
}