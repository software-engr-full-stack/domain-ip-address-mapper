terraform {
  required_version = ">= 1.1.5"

  cloud {
    organization = "software-engr-full-stack"
    workspaces {
      name = "aws-permanent"
    }
  }
}

module "config" { source = "../config" }

provider "aws" {
  region = module.config.region
}

resource "aws_eip" "terrainfra" {
  tags = {
    Name = module.config.eip_tag
  }
  vpc = true
}
