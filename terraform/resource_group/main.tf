resource "azurerm_resource_group" "main" {
  name     = "BugBlaster360RG"
  location = var.location

  tags = {
    "Name" = "BugBlaster360RG"
  }
}