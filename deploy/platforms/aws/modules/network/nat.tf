resource "aws_eip" "terrainfra_nat" {
  count = var.enable_nat ? 1 : 0

  # TODO: EIP may require IGW to exist prior to association. Use depends_on to set an explicit dependency on the IGW.
  #   https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/eip
  depends_on = [aws_internet_gateway.terrainfra]
  vpc = true
}

resource "aws_nat_gateway" "terrainfra" {
  count = var.enable_nat ? 1 : 0

  depends_on = [aws_internet_gateway.terrainfra]

  allocation_id = aws_eip.terrainfra_nat[count.index].id
  subnet_id     = aws_subnet.terrainfra_subnets_public["a"].id

  tags = {
    Name = var.name
  }
}
