resource "aws_internet_gateway" "terrainfra" {
  vpc_id = aws_vpc.terrainfra.id

  tags = {
    Name = var.name
  }
}
