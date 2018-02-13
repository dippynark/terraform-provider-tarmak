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

output "bastion_status" {
  value = "${data.tarmak_bastion_instance.bastion.status}"
}

data "tarmak_vault_cluster" "vault" {
  instances = ["ip-10-99-32-12.us-west-1.compute.internal",
		"ip-10-99-32-11.us-west-1.compute.internal",
		"ip-10-99-32-10.us-west-1.compute.internal"]
}

output "vault_status" {
  value = "${data.tarmak_vault_cluster.vault.status}"
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

output "vault_instance_role_etcd_status" {
  value = "${data.tarmak_vault_instance_role.etcd.status}"
}

output "vault_instance_role_master_status" {
  value = "${data.tarmak_vault_instance_role.master.status}"
}

output "vault_instance_role_worker_status" {
  value = "${data.tarmak_vault_instance_role.worker.status}"
}