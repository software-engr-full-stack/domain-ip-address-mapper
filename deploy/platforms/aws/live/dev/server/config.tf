variable "secrets" {
  type = object({
    aws_ssh_id_pub = string
    database = object({
      name = string
      user = string
      password = string
    })
  })
}

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

data "template_file" "user_data" {
  template = file("./cloud-init.yml")
  vars = {
    aws_ssh_id_pub = var.secrets.aws_ssh_id_pub
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-*20*-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

module "config" { source = "../../config" }

locals {
  name = "terrainfra-dev-server"

  instance = {
    key_name = "2022-02-01-terraform"

    # type = "t3.micro"
    type = "t2.micro"
    # type = "t2.nano"

    ami_image_id = data.aws_ami.ubuntu.id

    # # Red Hat flavored image
    # ami_image_id = "ami-01163e76c844a2129" # Amazon Linux 2 AMI (HVM) - Kernel 5.10, SSD Volume Type
  }

  ports = {
    ssh = 22
  }

  vpc_net_addr = "172.20"
}

locals {
  network = {
    cidr_block_vpc = "${local.vpc_net_addr}.0.0/16"
    subnets = {
       priv = {
        a = { az: "${module.config.region}a", cidr_block: "${local.vpc_net_addr}.10.0/24" }
        b = { az: "${module.config.region}b", cidr_block: "${local.vpc_net_addr}.20.0/24" }
      }

      public = {
        a = { az: "${module.config.region}a", cidr_block: "${local.vpc_net_addr}.40.240/28" }
        b = { az: "${module.config.region}b", cidr_block: "${local.vpc_net_addr}.80.240/28" }
      }
    }
  }
}
