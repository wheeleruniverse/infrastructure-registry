
data "google_iam_policy" "app" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

data "google_iam_policy" "web" {
  binding {
    role = "roles/storage.objectViewer"
    members = [
      "allUsers",
    ]
  }
}

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

resource "google_cloudbuild_trigger" "app-cicd" {
  description = "cloud build trigger for the app cicd"
  filename    = "cloudbuild.yml"
  name        = "${var.domain}-app-cicd"

  trigger_template {
    branch_name = "candidate"
    repo_name   = "wheelersadvice/cicd/app"
  }
}

resource "google_cloudbuild_trigger" "app-code" {
  description = "cloud build trigger for the app code"
  filename    = "cloudbuild.yml"
  name        = "${var.domain}-app-code"

  trigger_template {
    branch_name = "main"
    repo_name   = "github_wheelers-websites_cloudguruchallenge_21.08-app"
  }
}

resource "google_cloudbuild_trigger" "web-cicd" {
  description = "cloud build trigger for the web cicd"
  filename    = "cloudbuild.yml"
  name        = "${var.domain}-web-cicd"

  trigger_template {
    branch_name = "candidate"
    repo_name   = "wheelersadvice/cicd/web"
  }
}

resource "google_cloudbuild_trigger" "web-code" {
  description = "cloud build trigger for the web code"
  filename    = "cloudbuild.yml"
  name        = "${var.domain}-web-code"

  trigger_template {
    branch_name = "main"
    repo_name   = "github_wheelers-websites_cloudguruchallenge_21.08-web"
  }
}

resource "google_cloud_run_domain_mapping" "this" {
  location = google_cloud_run_service.this.location
  name     = "api.${var.fqdn}"

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

resource "google_cloud_run_service_iam_policy" "this" {
  location    = google_cloud_run_service.this.location
  policy_data = data.google_iam_policy.app.policy_data
  service     = google_cloud_run_service.this.name
}

resource "google_sourcerepo_repository" "app" {
  name = "${var.domain}/cicd/app"
}

resource "google_sourcerepo_repository" "web" {
  name = "${var.domain}/cicd/web"
}

resource "google_compute_backend_bucket" "this" {
  bucket_name = google_storage_bucket.this.name
  enable_cdn  = true
  name        = google_storage_bucket.this.name
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
      "${var.fqdn}.",
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

resource "google_dns_record_set" "api" {
  managed_zone = var.dns_zone
  name         = "api.${var.fqdn}."
  rrdatas      = ["ghs.googlehosted.com."]
  ttl          = 300
  type         = "CNAME"
}

resource "google_dns_record_set" "root" {
  managed_zone = var.dns_zone
  name         = "${var.fqdn}."
  rrdatas      = [google_compute_global_forwarding_rule.this.ip_address]
  ttl          = 300
  type         = "A"
}

resource "google_dns_record_set" "www" {
  managed_zone = var.dns_zone
  name         = "www.${var.fqdn}."
  rrdatas      = [google_compute_global_forwarding_rule.this.ip_address]
  ttl          = 300
  type         = "A"
}

resource "google_storage_bucket" "this" {
  name                        = var.domain
  force_destroy               = true
  labels                      = local.labels
  location                    = var.location_multi_region
  storage_class               = var.storage_class
  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
  }
}

resource "google_storage_bucket_iam_policy" "this" {
  bucket      = google_storage_bucket.this.name
  policy_data = data.google_iam_policy.web.policy_data
}
