resource "azurerm_storage_account" "main" {
  name                     = "bugblaster360sg"
  resource_group_name      = var.resource_group_name
  location                 = var.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    "Name" = "BugBlaster360SA"
  }
}