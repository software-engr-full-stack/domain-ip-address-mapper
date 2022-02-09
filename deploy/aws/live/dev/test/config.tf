variable "myip_mask" {
  type = number
  default = 32
  description = "The size of the bit mask (n.n.n.n/myip_mask)."
}

variable "allow_all" {
  type = bool
  description = "TODO: if false, only allow SSH and HTTP traffic from own IP. If true, TODO."
  default = false
}

variable "database" {
  type = object({
    name = string
    user = string
    password = string
  })
}

locals {
  name = "terrainfra-dev-test"
  region = "us-west-1"

  instance = {
    key_pair_name = "2022-02-01-terraform"
    type = "t2.nano"
    ami_image_id = "ami-01163e76c844a2129" # Amazon Linux 2 AMI (HVM) - Kernel 5.10, SSD Volume Type
  }

  ports = {
    http = 80
    ssh = 22
  }

  vpc_net_addr = "172.20"
}

locals {
  network = {
    test = {
      cidr_block_vpc = "${local.vpc_net_addr}.0.0/16"
      subnets = {
         priv = {
          a = { az: "${local.region}a", cidr_block: "${local.vpc_net_addr}.10.0/24" }
          b = { az: "${local.region}b", cidr_block: "${local.vpc_net_addr}.20.0/24" }
        }

        public = {
          a = { az: "${local.region}a", cidr_block: "${local.vpc_net_addr}.40.240/28" }
          b = { az: "${local.region}b", cidr_block: "${local.vpc_net_addr}.80.240/28" }
        }
      }
    }
  }
}
