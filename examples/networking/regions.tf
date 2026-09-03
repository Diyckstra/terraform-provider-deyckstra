# The provider is configured inside the region module, so the credentials
# read from c2rc.sh are passed to every module instance.
module "ru-msk" {
  source          = "./region"
  region          = "ru-msk"
  base_cidr_block = var.base_cidr_block

  access_key = "${local.dict_creds["C2_PROJECT"]}:${local.dict_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.dict_creds["EC2_SECRET_KEY"]
}

module "ru-spb" {
  source          = "./region"
  region          = "ru-spb"
  base_cidr_block = var.base_cidr_block

  access_key = "${local.dict_creds["C2_PROJECT"]}:${local.dict_creds["BASE_ACCESS_KEY"]}"
  secret_key = local.dict_creds["EC2_SECRET_KEY"]
}
