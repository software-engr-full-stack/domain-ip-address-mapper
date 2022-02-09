variable "name" { type = string }

variable "instance" {
  type = object({
    key_pair_name = string
    type = string
    ami_image_id = string
  })
}

variable "myip" {
  type = object({
    addr = string
    mask = number
  })
}

variable "ports" {
  type = object({
    ssh = number
    http = number
  })
}

variable "mod" {
  type = object({
    associate_public_ip_address = bool
    # user_data                   = string
  })
}

variable "data"   { type = any }

resource "aws_security_group" "terrainfra_svcs_instance" {
  vpc_id = var.data.vpc_id
  name = var.name

  ingress {
    from_port   = var.ports.http
    to_port     = var.ports.http
    protocol    = "tcp"
    cidr_blocks = ["${var.myip.addr}/${var.myip.mask}"]
  }

  ingress {
    from_port   = var.ports.ssh
    to_port     = var.ports.ssh
    protocol    = "tcp"
    cidr_blocks = ["${var.myip.addr}/${var.myip.mask}"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "terrainfra_svcs_instance" {
  key_name                    = var.instance.key_pair_name
  ami                         = var.instance.ami_image_id
  instance_type               = var.instance.type
  vpc_security_group_ids      = [aws_security_group.terrainfra_svcs_instance.id]
  subnet_id                   = var.data.subnets.public["a"].id
  associate_public_ip_address = var.mod.associate_public_ip_address
  # user_data                   = var.mod.user_data

  tags = {
    Name = var.name
  }
}
