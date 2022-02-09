terraform {
  required_version = ">= 1.1.5"

  backend "local" {
    path = "../../../state/live/dev/main.tfstate"
  }
}

data "external" "myip_addr" {
  count = var.allow_all ? file("... ERROR: TODO, allow_all traffic.") : 1
  program = ["${path.module}/../../../lib/own-ip.sh"]
}

locals {
  myip = {
    addr = data.external.myip_addr[0].result.myip_addr
    mask = var.myip_mask
  }
}

module "terrainfra_dev_network" {
  source = "../../../modules/network"
  name = local.name
  network_config = local.network.test
}

module "terrainfra_dev_bastion" {
  source = "../../../modules/security/bastion"

  config = {
    key_name = local.instance.key_pair_name
    name = local.name
    ssh_ingress_port = local.ports.ssh
    image_id = local.instance.ami_image_id
    instance_type = local.instance.type
  }

  myip = {
    addr = local.myip.addr
    mask = local.myip.mask
  }

  data = {
    vpc_id = module.terrainfra_dev_network.data.vpc_id
    subnet_id = module.terrainfra_dev_network.data.subnets.public["a"].id
  }
}

resource "aws_s3_bucket" "test-bucket" {
  bucket = "my-test-bucket-dingdong"
}
