
resource "google_app_engine_application" "this" {
  /*
   * Note:
   *
   * Two locations, which are called europe-west and us-central in App Engine commands and in the Google Cloud Console,
   * are called europe-west1 and us-central1, respectively, elsewhere in Google documentation.
   *
   * https://cloud.google.com/appengine/docs/locations
   */
  database_type = "CLOUD_FIRESTORE"
  location_id   = "US-CENTRAL1" == var.location_region ? "us-central" : lower(var.location_region)
}

resource "google_cloud_run_domain_mapping" "this" {
  location = google_cloud_run_service.this.location
  name     = "api.gcp.${var.fqdn}"

  metadata {
    namespace = var.project
  }
  spec {
    route_name = google_cloud_run_service.this.name
  }
}

resource "google_cloud_run_service" "this" {
  autogenerate_revision_name = true
  location                   = lower(var.location_region)
  name                       = var.domain

  template {
    metadata {
      labels = local.labels
    }
    spec {
      containers {
        image = var.image

        env {
          name  = "GCP_CREDENTIALS"
          value = var.service_account_json
        }
      }
    }
  }
  traffic {
    latest_revision = true
    percent         = 100
  }
}

resource "google_compute_backend_bucket" "this" {
  bucket_name = google_storage_bucket.web.name
  enable_cdn  = true
  name        = google_storage_bucket.web.name
}

resource "google_compute_global_forwarding_rule" "this" {
  name       = var.domain
  target     = google_compute_target_https_proxy.this.id
  port_range = "443"
}

resource "google_compute_managed_ssl_certificate" "this" {
  name = var.domain

  managed {
    domains = [
      "gcp.${var.fqdn}.",
      "www.${var.fqdn}."
    ]
  }
}

resource "google_compute_target_https_proxy" "this" {
  name             = var.domain
  ssl_certificates = [google_compute_managed_ssl_certificate.this.id]
  url_map          = google_compute_url_map.this.id
}

resource "google_compute_url_map" "this" {
  default_service = google_compute_backend_bucket.this.id
  name            = var.domain
}

resource "google_container_registry" "this" {
  location = var.location_multi_region
}

resource "google_storage_bucket" "app" {
  name                        = "${var.domain}-app"
  force_destroy               = true
  labels                      = local.labels
  location                    = var.location_multi_region
  storage_class               = var.storage_class
  uniform_bucket_level_access = true
}

resource "google_storage_bucket" "iac" {
  name                        = "${var.domain}-iac"
  force_destroy               = true
  labels                      = local.labels
  location                    = var.location_multi_region
  storage_class               = var.storage_class
  uniform_bucket_level_access = true
}

resource "google_storage_bucket" "web" {
  name                        = "${var.domain}-web"
  force_destroy               = true
  labels                      = local.labels
  location                    = var.location_multi_region
  storage_class               = var.storage_class
  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
  }
}
