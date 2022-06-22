
resource "google_redis_instance" "redis" {
  name           = "config-api-memory-cache"
  tier           = "STANDARD_HA"
  memory_size_gb = 1

  location_id             = var.zone

  authorized_network = data.google_compute_network.redis-network.id

  redis_version     = "REDIS_6_X"
  display_name      = "Config API Cache Instance"


}

output "redis_host" {
  value = google_redis_instance.redis.host
}

output "redis_port" {
  value = google_redis_instance.redis.port
}

output "redis_current_location_id" {
  value = google_redis_instance.redis.current_location_id
}

data "google_compute_network" "redis-network" {
  name = "redis-test-network"
}