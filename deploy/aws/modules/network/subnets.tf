resource "aws_subnet" "terrainfra_subnets_priv" {
  vpc_id = aws_vpc.terrainfra.id

  for_each = var.network_config.subnets.priv

  cidr_block = lookup(each.value, "cidr_block")
  availability_zone = lookup(each.value, "az")

  tags = {
    Name = "${var.name}-priv-${each.key}"
  }
}

resource "aws_subnet" "terrainfra_subnets_public" {
  vpc_id = aws_vpc.terrainfra.id

  for_each = var.network_config.subnets.public

  cidr_block = lookup(each.value, "cidr_block")
  availability_zone = lookup(each.value, "az")
  map_public_ip_on_launch = true

  tags = {
    Name = "${var.name}-public-${each.key}"
  }
}
