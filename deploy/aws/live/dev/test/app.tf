# module "terrainfra_dev_db" {
#   source = "../../../modules/services/db/postgres"

#   name = local.name
#   data = module.terrainfra_dev_network.data

#   mod = {
#     name = var.database.name
#     username = var.database.user
#     password = var.database.password

#     postgres_port = 5432
#   }
# }

module "terrainfra_dev_test_instance" {
  source = "../../../modules/services/instance"
  name = local.name
  instance = local.instance
  myip = local.myip
  ports = local.ports
  data = module.terrainfra_dev_network.data
  mod = {
    associate_public_ip_address = true

    # user_data = templatefile(
    #   "${path.module}/user_data.tmpl.sh", {
    #     tf_deploy_user = var.deploy.user

    #     tf_ssh_ingress_port = local.config.port.ssh

    #     tf_http_port = 8080
    #     # tf_http_port = local.config.port.http

    #     tf_do_db_stuff = var.deploy.db.do_db_stuff

    #     # endpoint - The connection endpoint in address:port format.
    #     # url: {{envOr "DATABASE_URL" "postgres://postgres:postgres@127.0.0.1:5432/awsexample_production?sslmode=disable"}}
    #     # "postgres://YourUserName:YourPassword@YourHost:5432/YourDatabase";
    #     tf_database_url = join("", [
    #       "postgres://",
    #       var.deploy.user,
    #       ":",
    #       var.deploy.db.password,
    #       "@",
    #       module.terrainfra_dev_db.data.endpoint,
    #       "/",
    #       var.deploy.db.name,
    #       "?sslmode=disable"
    #     ])
    #   }
    # )
  }
}
