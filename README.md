# terraform-provider-tarmak

terraform-provider-tarmak is a Terraform provider built specifically for use within [Tarmak](https://github.com/jetstack/tarmak).

## Examples

```
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
```