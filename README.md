# terraform-provider-tarmak

terraform-provider-tarmak is a Terraform provider built specifically for use within [Tarmak](https://github.com/jetstack/tarmak). 

## Examples

```
resource "tarmak_tunnel" "primary" {
  bind_address = "127.0.0.1"
  destination_address = "216.58.206.110"
  destination_port = 80
  ssh_config_path = "/Users/root/.ssh/config"
  ssh_config_host = "server"
}
```