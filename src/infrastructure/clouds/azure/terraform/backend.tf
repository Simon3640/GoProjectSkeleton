terraform {
  backend "azurerm" {
    resource_group_name  = "tfstate-rg"
    storage_account_name = "tfstateaccountsg001"
    container_name       = "tfstate"
    key                  = "infra.tfstate"
  }
}
