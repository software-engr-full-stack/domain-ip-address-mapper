variable "name" { type = string }

variable "data"   { type = any }

variable "mod" {
  type = object({
    postgres_port = number
    name = string
    username = string
    password = string
  })
}

resource "aws_security_group" "terrainfra_svcs_db_postgres" {
  vpc_id = var.data.vpc_id
  name = "${var.name}-db-postgres"

  ingress {
    from_port   = var.mod.postgres_port
    to_port     = var.mod.postgres_port
    protocol    = "tcp"
    cidr_blocks = concat(
      [for key, snet in var.data.subnets.priv : snet.cidr_block],
      [for key, snet in var.data.subnets.public : snet.cidr_block]
    )
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_db_subnet_group" "terrainfra_svcs_db_postgres" {
  name       = var.name
  subnet_ids = [for key, snet in var.data.subnets.priv : snet.id]

  tags = {
    Name = var.name
  }
}

resource "aws_db_instance" "terrainfra_svcs_db_postgres" {
  # See docs.aws.amazon.com/AmazonRDS/latest/APIReference/API_CreateDBInstance.html for valid values
  db_subnet_group_name = aws_db_subnet_group.terrainfra_svcs_db_postgres.name
  vpc_security_group_ids = [aws_security_group.terrainfra_svcs_db_postgres.id]
  # publicly_accessible = false # ... default, present for doc.

  skip_final_snapshot = true # Probably set to true in real production.
  # final_snapshot_identifier = "???" # Must be set if skip_final_snapshot is false

  # storage_encrypted = "???" # Probably set to true in real production.
  allocated_storage = 20 # ... GB, minimum
  max_allocated_storage = 0  # ... to disable storage autoscaling
  storage_type = "gp2"
  engine = "postgres"
  # engine_version = "???"

  # ... 2022 02 14, $ 0.021/hour
  instance_class = "db.t4g.micro"

  # deletion_protection = "???" # Probably set to true in real production.
  # enabled_cloudwatch_logs_exports = "???" # Probably assign value in real production.

  # TODO: refactor
  db_name = var.mod.name
  username = var.mod.username
  password = var.mod.password
  # parameter_group_name = "???" # Should be set, ideally.

  auto_minor_version_upgrade = true
  # copy_tags_to_snapshot = ??? # Probably set to true in real production.

  tags = {
    Name = "terrainfra"
  }

  # See Terraform and AWS docs for rest of options (concerning performance insights, etc.)
}

output "output" {
  value = aws_db_instance.terrainfra_svcs_db_postgres
}
