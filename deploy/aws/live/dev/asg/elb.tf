resource "aws_elb" "terrainfra" {
  count           = var.enable_elb ? 1 : 0

  security_groups = [aws_security_group.terrainfra_elb[count.index].id]
  name            = var.name

  subnets = [for key, snet in aws_subnet.terrainfra_subnets_public : snet.id]

  health_check {
    target              = "HTTP:${var.elb_ports.instance}/"
    interval            = 30
    timeout             = 3
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }

  listener {
    lb_port           = var.elb_ports.listener
    lb_protocol       = "http"
    instance_port     = var.elb_ports.instance
    instance_protocol = "http"
  }
}

resource "aws_security_group" "terrainfra_elb" {
  count = var.enable_elb ? 1 : 0

  vpc_id = aws_vpc.terrainfra.id
  name = "${var.name}-elb"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = var.elb_ports.listener
    to_port     = var.elb_ports.listener
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
