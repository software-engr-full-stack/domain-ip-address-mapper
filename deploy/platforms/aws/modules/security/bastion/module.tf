variable "config" {
  type = object({
    key_name = string
    name = string
    image_id = string
    instance_type = string
    ssh_ingress_port = number
    user_data = string
  })
}

variable "myip" {
  type = object({
    addr = string
    mask = number
  })
}

variable "data" {
  type = object({
    vpc_id = string
    subnet_id = string
  })
}

resource "aws_instance" "terrainfra_bastion" {
  key_name               = var.config.key_name
  ami                    = var.config.image_id
  instance_type          = var.config.instance_type
  vpc_security_group_ids = [aws_security_group.terrainfra_bastion.id]
  subnet_id              = var.data.subnet_id

  user_data = var.config.user_data

  tags = {
    Name = "${var.config.name}-bastion"
  }
}

resource "aws_security_group" "terrainfra_bastion" {
  vpc_id = var.data.vpc_id
  name = "${var.config.name}-bastion"

  ingress {
    from_port   = var.config.ssh_ingress_port
    to_port     = var.config.ssh_ingress_port
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

output "bastion" {
  value = aws_instance.terrainfra_bastion.public_ip
}
