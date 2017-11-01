module "tarmak_vault" {
  source = "./tarmak_vault"
}

resource "null_resource" "worker_token_to_file" {
    provisioner "local-exec" {
        command = "echo ${module.tarmak_vault.init_token} >> token_test.txt"
    }
}
