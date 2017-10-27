resource "tarmak_vault_init_token" "worker" {
    environment = "josh"
    cluster = "cluster"
}

resource "null_resource" "worker_token_to_file" {
    provisioner "local-exec" {
        command = "echo ${tarmak_vault_init_token.worker.init_token} >> token_test.txt"
    }
}
