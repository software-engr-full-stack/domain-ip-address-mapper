variable "name" { type = string }

variable "instance" {
  type = object({
    key_name = string
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
  })
}

variable "user_data" {
  type = string
  default = ""
}

variable "eip" {
  type = object({
    do_associate  = bool
    allocation_id = string
  })

  default = {
    do_associate = false
    allocation_id = 0
  }
}

variable "associate_public_ip_address" {
  type = bool
  default = false
}

variable "data"   { type = any }

resource "aws_security_group" "terrainfra_svcs_instance" {
  vpc_id = var.data.vpc_id
  name = var.name

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
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

resource "aws_eip_association" "eip_assoc" {
  count = var.eip.do_associate ? 1 : 0
  instance_id   = aws_instance.terrainfra_svcs_instance.id
  allocation_id = var.eip.allocation_id
}

resource "aws_instance" "terrainfra_svcs_instance" {
  key_name                    = var.instance.key_name
  ami                         = var.instance.ami_image_id
  instance_type               = var.instance.type
  vpc_security_group_ids      = [aws_security_group.terrainfra_svcs_instance.id]
  subnet_id                   = var.data.subnets.public["a"].id

  # You can probably skip this in the caller module if using elastic IP (EIP).
  # But it might always destroy/create the instance depending on the default value.
  # If using EIP with public address, this will always become true.
  associate_public_ip_address = var.eip.do_associate || var.associate_public_ip_address

  user_data                   = var.user_data

  tags = {
    Name = var.name
  }
}

output "output" {
  value = aws_instance.terrainfra_svcs_instance
}
