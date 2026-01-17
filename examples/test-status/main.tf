terraform {
  required_providers {
    sonarr = {
      source  = "registry.terraform.io/oleksii-kalinin/sonarr"
      version = "0.0.1"
    }
  }
}

provider "sonarr" {
  url     = ""         // Sonarr URL
  api_key = "" // Sonarr API Key
}

data "sonarr_system_status" "sonarr" {}

data "sonarr_series_lookup" "lost" {
  term = "lost"
}

resource "sonarr_series" "lost" {
  tvdb_id         = data.sonarr_series_lookup.lost.tvdb_id
  path            = "/media/series"
  quality_profile = 1
  title           = data.sonarr_series_lookup.lost.title
  monitored       = true
}

output "lost_info" {
  value = data.sonarr_series_lookup.lost
}

output "sonarr_info" {
  value = data.sonarr_system_status.sonarr
}
