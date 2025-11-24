# Azure Redis Cache
resource "azurerm_redis_cache" "redis" {
  name                = "${local.name_prefix}-redis"
  location            = var.location
  resource_group_name = azurerm_resource_group.rg.name
  sku_name            = "Basic"
  capacity            = 1
  family              = "C"
  redis_version       = "6"
}
