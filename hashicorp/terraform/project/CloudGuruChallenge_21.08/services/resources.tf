
resource "google_billing_budget" "this" {
  billing_account = var.billing_account
  display_name    = var.domain

  all_updates_rule {
    disable_default_iam_recipients = true
    monitoring_notification_channels = [
      google_monitoring_notification_channel.this.id,
    ]
  }
  amount {
    specified_amount {
      currency_code = "USD"
      units         = var.budget
    }
  }
  budget_filter {
    credit_types_treatment = "EXCLUDE_ALL_CREDITS"
    projects               = ["projects/${var.project}"]
  }
  threshold_rules {
    threshold_percent = 0.25
  }
  threshold_rules {
    threshold_percent = 0.50
  }
  threshold_rules {
    threshold_percent = 0.75
  }
  threshold_rules {
    threshold_percent = 1.00
  }
}

resource "google_monitoring_notification_channel" "this" {
  display_name = var.owner
  labels = {
    email_address = var.owner
  }
  type = "email"
}

resource "google_project_iam_member" "build" {
  count  = length(local.build_roles)
  member = "serviceAccount:${var.project_number}@cloudbuild.gserviceaccount.com"
  role   = "roles/${local.build_roles[count.index]}"
}

resource "google_project_iam_member" "core" {
  count  = length(local.core_roles)
  member = "serviceAccount:${google_service_account.this.email}"
  role   = "roles/${local.core_roles[count.index]}"
}

resource "google_project_service" "this" {
  count                      = length(local.services)
  disable_dependent_services = true
  disable_on_destroy         = true
  service                    = "${local.services[count.index]}.googleapis.com"

  timeouts {
    create = local.timeout
    update = local.timeout
  }
}

resource "google_service_account" "this" {
  account_id   = "${var.domain}-svc"
  description  = "service account for ${var.domain}"
  display_name = "${var.domain}-svc"
}

resource "google_service_account_key" "this" {
  service_account_id = google_service_account.this.name
}
