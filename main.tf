variable "cluster" {
    type = "string"
    default = "cluster"
}

variable "environment" {
    type = "string"
    default = "luke"
}

provider "tarmak" {
  environment = "${var.environment}"
  cluster     = "${var.cluster}"
}

resource "tarmak_bastion_instance" "bastion" {  
}

output "bastion_status" {
  value = "${tarmak_bastion_instance.bastion.status}"
}


/*

resource "tarmak_tunnel" "test" {
  bind_address = "127.0.0.1"
  #destination_address = "google.com"
  destination_address = "216.58.206.110"
  destination_port = 80
  ssh_config_path = "/Users/luke/.ssh/config"
  ssh_config_host = "home"
}

output "bind_port" {
  value = "${tarmak_tunnel.test.bind_port}"
}


variable "role" {
    type = "string"
    default = "worker"
}

resource "tarmak_vault_init_token" "worker" {
  environment = "${var.environment}"
  cluster     = "${var.cluster}"
  role        = "${var.role}"
}

resource "null_resource" "worker_token_to_file" {
  provisioner "local-exec" {
    command = "echo ${tarmak_vault_init_token.worker.init_token} >> token_test.txt"
  }
}
*/
