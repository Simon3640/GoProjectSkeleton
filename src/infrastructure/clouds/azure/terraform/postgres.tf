resource "azurerm_postgresql_flexible_server" "postgres" {
  name                = "${local.name_prefix}-pg"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  version             = "15"

  administrator_login    = var.db_user
  administrator_password = azurerm_key_vault_secret.db_password.value

  sku_name = "B_Standard_B1ms"  # ✔ Más barato compatible con free tier

  storage_mb = 32768            # 32 GB mínimo

  backup_retention_days  = 7

  tags = {
    environment = var.environment
  }

  # Ignorar cambios en zone para evitar errores al importar recursos existentes
  lifecycle {
    ignore_changes = [zone]
  }

  depends_on = [azurerm_key_vault_secret.db_password]
}

# Crear la base de datos
resource "azurerm_postgresql_flexible_server_database" "main" {
  name      = var.db_name
  server_id = azurerm_postgresql_flexible_server.postgres.id
  charset   = "UTF8"
  collation = "en_US.utf8"

  depends_on = [azurerm_postgresql_flexible_server.postgres]
}

# Permite todo Internet temporalmente
# Lo ideal: agregar sólo la IP de tu Function App (si está fija)
resource "azurerm_postgresql_flexible_server_firewall_rule" "allow_all" {
  name                = "allow-all"
  server_id           = azurerm_postgresql_flexible_server.postgres.id
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "255.255.255.255"
}
