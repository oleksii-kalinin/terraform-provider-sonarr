terraform {
  required_providers {
    sonarr = {
      source = "registry.terraform.io/oleksii-kalinin/sonarr"
      version = "0.0.1"
    }
  }
}

provider "sonarr" {
  url = "" // Sonarr URL
  api_key = "" // Sonarr API Key
}

data "sonarr_system_status" "sonarr" {}

resource "sonarr_series" "mr-robot" {
  tvdb_id = "289590"
  path = "/media/series"
  quality_profile = 1
  title = "Mr. Robot"
  monitored = true
  add_options = {
    monitor = "all"
  }
}

output "sonarr_info" {
  value = data.sonarr_system_status.sonarr
}