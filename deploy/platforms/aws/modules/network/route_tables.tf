resource "aws_route_table" "terrainfra_priv" {
  vpc_id = aws_vpc.terrainfra.id

  tags = {
    Name = "${var.name}-priv"
  }
}

resource "aws_route" "terrainfra_priv_nat_gateway" {
  count = var.enable_nat ? 1 : 0

  route_table_id            = aws_route_table.terrainfra_priv.id
  destination_cidr_block    = "0.0.0.0/0"
  nat_gateway_id = aws_nat_gateway.terrainfra[count.index].id
}


resource "aws_route_table_association" "terrainfra_priv" {
  route_table_id = aws_route_table.terrainfra_priv.id

  for_each  = aws_subnet.terrainfra_subnets_priv
  subnet_id = each.value.id
}

resource "aws_route_table" "terrainfra_public" {
  vpc_id = aws_vpc.terrainfra.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.terrainfra.id
  }

  tags = {
    Name = "${var.name}-public"
  }
}

resource "aws_route_table_association" "terrainfra_public" {
  route_table_id = aws_route_table.terrainfra_public.id

  for_each  = aws_subnet.terrainfra_subnets_public
  subnet_id = each.value.id
}
