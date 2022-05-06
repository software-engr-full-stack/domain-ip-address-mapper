module "terrainfra_dev_db" {
  source = "../../../modules/services/db/postgres"

  name = local.name
  data = module.terrainfra_dev_network.data

  mod = {
    name = var.secrets.database.name
    username = var.secrets.database.user
    password = var.secrets.database.password

    postgres_port = 5432
  }
}

data "aws_eip" "by_tag" {
  filter {
    name   = "tag:Name"
    values = [module.config.eip_tag]
  }
}

module "terrainfra_dev_server_instance" {
  source = "../../../modules/services/instance"
  name = local.name
  instance = local.instance
  myip = local.myip
  ports = local.ports
  data = module.terrainfra_dev_network.data

  eip = {
    do_associate = true
    allocation_id = data.aws_eip.by_tag.id
  }
  user_data = data.template_file.user_data.rendered
}
