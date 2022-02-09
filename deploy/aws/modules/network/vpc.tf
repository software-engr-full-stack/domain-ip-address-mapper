resource "aws_vpc" "terrainfra" {
  cidr_block = var.network_config.cidr_block_vpc

  enable_dns_hostnames = var.enable_dns_hostnames
  enable_dns_support = var.enable_dns_support
  assign_generated_ipv6_cidr_block = true

  tags = {
    Name = var.name
  }
}
