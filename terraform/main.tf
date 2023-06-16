module "resource_group" {
  source   = "./resource_group"
  location = var.location
}

module "storage_account" {
  source              = "./storage_account"
  location            = var.location
  resource_group_name = module.resource_group.resource_group_name
}