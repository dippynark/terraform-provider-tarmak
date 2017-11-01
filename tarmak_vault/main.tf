resource "tarmak_vault_init_token" "worker" {
    environment = "josh"
    cluster = "cluster"
    role = "worker"
}

output "init_token" {
  value = "${tarmak_vault_init_token.worker.init_token}"
}
