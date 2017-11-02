variable "cluster" {
    type = "string"
    default = "cluster"
}

variable "environment" {
    type = "string"
    default = "josh"
}

variable "role" {
    type = "string"
    default = "worker"
}

provider "tarmak" {
  environment = "${var.environment}"
  cluster     = "${var.cluster}"
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
