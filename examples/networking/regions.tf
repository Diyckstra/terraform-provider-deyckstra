module "ru-msk" {
  source          = "./region"
  region          = "ru-msk"
  base_cidr_block = var.base_cidr_block
}

module "ru-spb" {
  source          = "./region"
  region          = "ru-spb"
  base_cidr_block = var.base_cidr_block
}
