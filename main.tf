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

data "tarmak_bastion_instance" "bastion" {  
  hostname = "bastion"
  username = "root"
}

data "tarmak_vault_cluster" "vault" {
  instances = ["ip-10-99-32-12.us-west-1.compute.internal",
		"ip-10-99-32-11.us-west-1.compute.internal",
		"ip-10-99-32-10.us-west-1.compute.internal"]
}

data "tarmak_vault_instance_role" "etcd" {
  vault_cluster_name = "vault"
  role_name = "etcd" 
}

data "tarmak_vault_instance_role" "master" {
  vault_cluster_name = "vault"
  role_name = "master" 
}

data "tarmak_vault_instance_role" "worker" {
  vault_cluster_name = "vault"
  role_name = "worker" 
}

output "bastion_instance_status" {
  value = "${data.tarmak_bastion_instance.bastion.status}"
}

output "vault_cluster_status" {
  value = "${data.tarmak_vault_cluster.vault.status}"
}

output "vault_instance_role_etcd_init_token" {
  value = "${data.tarmak_vault_instance_role.etcd.init_token}"
}

output "vault_instance_role_master_init_token" {
  value = "${data.tarmak_vault_instance_role.master.init_token}"
}

output "vault_instance_role_worker_init_token" {
  value = "${data.tarmak_vault_instance_role.worker.init_token}"
}