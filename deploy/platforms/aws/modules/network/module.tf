variable "name" {
  type = string
}

variable "network_config" {
  type = object({
    cidr_block_vpc = string
    subnets = map(any)
  })
}

variable "enable_nat" {
  type = bool
  default = true
}

variable "enable_dns_hostnames" {
  type = bool
  default = false
}

variable "enable_dns_support" {
  type = bool
  # ... 2022 03 13: I have no idea why I originally put default = false here.
  default = true
  # default = false
}

output "data" {
  value = {
    vpc_id = aws_vpc.terrainfra.id
    subnets = {
      priv: aws_subnet.terrainfra_subnets_priv,
      public: aws_subnet.terrainfra_subnets_public
    }
  }
}
