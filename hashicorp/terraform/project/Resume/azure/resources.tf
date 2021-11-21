
resource "azurerm_application_insights" "this" {
  application_type    = "java"
  location            = azurerm_resource_group.this.location
  name                = var.domain
  resource_group_name = azurerm_resource_group.this.name
  retention_in_days   = 30
  sampling_percentage = 1
  tags                = local.tags
}

resource "azurerm_app_service_plan" "this" {
  kind                = "FunctionApp"
  location            = azurerm_resource_group.this.location
  name                = var.domain
  resource_group_name = azurerm_resource_group.this.name
  tags                = local.tags

  sku {
    size = "Y1"
    tier = "Dynamic"
  }
}

resource "azurerm_cdn_endpoint" "this" {
  content_types_to_compress = [
    "application/javascript",
    "application/json",
    "application/x-javascript",
    "application/xml",
    "text/css",
    "text/html",
    "text/javascript",
    "text/plain"
  ]
  is_compression_enabled = true
  location               = azurerm_resource_group.this.location
  name                   = var.domain
  origin_host_header     = "${azurerm_storage_account.web.name}.z20.web.core.windows.net"
  profile_name           = azurerm_cdn_profile.this.name
  resource_group_name    = azurerm_resource_group.this.name
  tags                   = local.tags

  delivery_rule {
    name  = "HttpToHttpsRedirect"
    order = 1

    request_scheme_condition {
      match_values = ["HTTP"]
      operator     = "Equal"
    }
    url_redirect_action {
      protocol      = "Https"
      redirect_type = "PermanentRedirect"
    }
  }
  delivery_rule {
    name  = "SetCacheExpiration"
    order = 2

    cache_expiration_action {
      behavior = "SetIfMissing"
      duration = "30.00:00:00"
    }
    url_file_extension_condition {
      match_values = [
        "css",
        "ico",
        "js",
      ]
      negate_condition = false
      operator         = "Equal"
      transforms       = []
    }
  }
  global_delivery_rule {
    modify_response_header_action {
      action = "Append"
      name   = "Strict-Transport-Security"
      value  = "max-age=63072000; includeSubDomains; preload"
    }
    modify_response_header_action {
      action = "Append"
      name   = "X-Content-Type-Options"
      value  = "nosniff"
    }
    modify_response_header_action {
      action = "Append"
      name   = "X-Frame-Options"
      value  = "SAMEORIGIN"
    }
    modify_response_header_action {
      action = "Append"
      name   = "X-XSS-Protection"
      value  = "1; mode=block"
    }
  }
  origin {
    host_name = "${azurerm_storage_account.web.name}.z20.web.core.windows.net"
    name      = var.domain
  }
}

resource "azurerm_cdn_profile" "this" {
  location            = azurerm_resource_group.this.location
  name                = var.domain
  resource_group_name = azurerm_resource_group.this.name
  sku                 = "Standard_Microsoft"
  tags                = local.tags
}

resource "azurerm_cosmosdb_account" "this" {
  enable_automatic_failover = false
  enable_free_tier          = true
  kind                      = "GlobalDocumentDB"
  location                  = azurerm_resource_group.this.location
  name                      = var.domain
  offer_type                = "Standard"
  resource_group_name       = azurerm_resource_group.this.name
  tags                      = local.tags

  capabilities {
    name = "EnableServerless"
  }
  consistency_policy {
    consistency_level = "Eventual"
  }
  geo_location {
    failover_priority = 0
    location          = azurerm_resource_group.this.location
  }
}

resource "azurerm_cosmosdb_sql_container" "this" {
  account_name        = azurerm_cosmosdb_account.this.name
  count               = length(var.tables)
  database_name       = azurerm_cosmosdb_sql_database.this.name
  name                = var.tables[count.index]
  partition_key_path  = "/id"
  resource_group_name = azurerm_resource_group.this.name
}

resource "azurerm_cosmosdb_sql_database" "this" {
  account_name        = azurerm_cosmosdb_account.this.name
  name                = var.project
  resource_group_name = azurerm_resource_group.this.name
}

resource "azurerm_function_app" "this" {
  app_service_plan_id = azurerm_app_service_plan.this.id
  app_settings = {
    APPINSIGHTS_INSTRUMENTATIONKEY = azurerm_application_insights.this.instrumentation_key
    COSMOS_AUTH                    = azurerm_cosmosdb_account.this.primary_master_key
    COSMOS_HOST                    = azurerm_cosmosdb_account.this.endpoint
    COSMOS_NAME                    = azurerm_cosmosdb_sql_database.this.name
    ENABLE_ORYX_BUILD              = false
    FUNCTIONS_WORKER_RUNTIME       = "java"
    FUNCTION_APP_EDIT_MODE         = "readonly"
    FUNCTIONS_EXTENSION_VERSION    = "~3"
    SCM_DO_BUILD_DURING_DEPLOYMENT = false
    WEBSITE_RUN_FROM_PACKAGE       = "1"
  }
  https_only                 = true
  location                   = azurerm_resource_group.this.location
  name                       = var.domain
  resource_group_name        = azurerm_resource_group.this.name
  storage_account_access_key = azurerm_storage_account.app.primary_access_key
  storage_account_name       = azurerm_storage_account.app.name
  tags                       = local.tags
  version                    = "~3"

  site_config {
    java_version = "11"

    cors {
      allowed_origins     = ["*"]
      support_credentials = false
    }
  }
}

resource "azurerm_resource_group" "this" {
  location = "East US 2"
  name     = var.domain
  tags     = local.tags
}

resource "azurerm_storage_account" "app" {
  account_replication_type = var.bucket_replication
  account_tier             = var.bucket_tier
  allow_blob_public_access = false
  location                 = azurerm_resource_group.this.location
  min_tls_version          = var.tls_version
  name                     = "${var.domain}app"
  resource_group_name      = azurerm_resource_group.this.name
  tags                     = local.tags

  network_rules {
    bypass         = ["AzureServices"]
    default_action = "Allow"
  }
}

resource "azurerm_storage_account" "iac" {
  account_replication_type = var.bucket_replication
  account_tier             = var.bucket_tier
  allow_blob_public_access = true
  location                 = azurerm_resource_group.this.location
  min_tls_version          = var.tls_version
  name                     = "${var.domain}iac"
  resource_group_name      = azurerm_resource_group.this.name
  tags                     = local.tags

  network_rules {
    bypass         = ["AzureServices"]
    default_action = "Allow"
  }
}

resource "azurerm_storage_account" "web" {
  account_replication_type = var.bucket_replication
  account_tier             = var.bucket_tier
  allow_blob_public_access = false
  location                 = azurerm_resource_group.this.location
  min_tls_version          = var.tls_version
  name                     = "${var.domain}web"
  resource_group_name      = azurerm_resource_group.this.name
  tags                     = local.tags

  network_rules {
    bypass         = ["AzureServices"]
    default_action = "Allow"
  }
  static_website {
    index_document = "index.html"
  }
}

resource "azurerm_storage_container" "this" {
  container_access_type = "container"
  name                  = "web"
  storage_account_name  = azurerm_storage_account.iac.name
}
