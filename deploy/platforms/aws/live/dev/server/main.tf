terraform {
  required_version = ">= 1.1.5"

  cloud {
    organization = "software-engr-full-stack"
    workspaces {
      name = "aws-server"
    }
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
  network_config = local.network
  enable_nat = false
}

output "devices" {
  value = {
    eip = {
      public_ip = data.aws_eip.by_tag.public_ip
      public_dns = data.aws_eip.by_tag.public_dns
    }
    server = module.terrainfra_dev_server_instance.output.public_ip
    database = module.terrainfra_dev_db.output.endpoint
  }
}
