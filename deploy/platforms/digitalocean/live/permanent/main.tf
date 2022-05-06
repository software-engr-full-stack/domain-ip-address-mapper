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
      name = "digitalocean-permanent"
    }
  }
}

variable "digitalocean_secret_token" {
  type = string
}

provider "digitalocean" {
  token = var.digitalocean_secret_token
}

variable "region" {
  type = string
}

# Unused, present to silence warnings
variable "floating_ip" {
  type = string
}

resource "digitalocean_floating_ip" "cloudsandbox" {
  region = var.region
}

output "floating_ip" {
  value = digitalocean_floating_ip.cloudsandbox
}
