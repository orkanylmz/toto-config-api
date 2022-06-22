
module cloud_run_config_api_http {
  source = "./service"

  project    = var.project
  location   = var.region
  dependency = null_resource.init_docker_images


  name     = "skuconfig"
  protocol = "http"
  auth     = false

  envs = [
    {
      name = "DB_CONN_STRING"
      value =  google_sql_database_instance.postgresql.connection_name
    },
    {
      name = "POSTGRES_USER"
      value = var.db_user_name
    },
    {
      name = "POSTGRES_DB"
      value = var.db_name
    },
    {
      name = "POSTGRES_PASSWORD"
      value = var.db_user_password
    },
    {
      name = "DB_DRIVER"
      value = "cloudsqlpostgres"
    },
    {
      name = "REDIS_HOST"
      value = "${google_redis_instance.redis.host}:${google_redis_instance.redis.port}"
    },
  ]

}

