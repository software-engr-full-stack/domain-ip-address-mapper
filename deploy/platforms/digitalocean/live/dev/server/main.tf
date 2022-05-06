terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }

  cloud {
    organization = "software-engr-full-stack"
    workspaces {
      name = "digitalocean-server"
    }
  }
}

variable "digitalocean_secret_token" {}

provider "digitalocean" {
  token = var.digitalocean_secret_token
}

# ... Command": ssh-keygen -E md5 -lf ~/.ssh/id_ed25519.pub | awk '{print $2}' | cut -f 2- -d ':'
variable "ssh_fingerprints" {
  type = list(string)
  description = "The fingerprints of the SSH pub keys used to connect to droplets."
}

# TODO: derive this from the fingerprints or derive the fingerprints from these keys.
variable "ssh_authorized_keys" {
  type = list(string)
}

variable "region" {
  type = string
}

variable "floating_ip" {
  type = string
}

data "digitalocean_floating_ip" "cloudsandbox" {
  ip_address = var.floating_ip
}

resource "digitalocean_vpc" "vpc" {
  name     = "vpc"
  region   = var.region
  ip_range = local.network.vpc_multi.ip_range
}

data "template_file" "user_data" {
  template = file("./cloud-init.yml")
  vars = {
    ssh_authorized_key0 = var.ssh_authorized_keys[0]
    ssh_authorized_key1 = var.ssh_authorized_keys[1]
  }
}

resource "digitalocean_droplet" "server" {
  ssh_keys = var.ssh_fingerprints
  size = "s-1vcpu-1gb"
  image = "ubuntu-20-04-x64"
  region = var.region
  vpc_uuid = digitalocean_vpc.vpc.id
  user_data = data.template_file.user_data.rendered

  name = "server"
}

resource "digitalocean_floating_ip_assignment" "fip_assignment" {
  ip_address = data.digitalocean_floating_ip.cloudsandbox.ip_address
  droplet_id = digitalocean_droplet.server.id
}

output "server" {
  value = digitalocean_droplet.server
}
