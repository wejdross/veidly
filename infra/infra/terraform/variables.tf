variable "hcloud_token" {
  sensitive = true
}
variable "hetzner_dns_token" {
  sensitive = true
}
variable "veidly_infra_machines" {
  type        = map(any)
  description = "values shared across all machines"
  default = {
    data_volume_size = 20
    image            = "ubuntu-20.04"
    location         = "hel1"
    server_type      = "cpx11"
  }
}

variable "vms" {
  type        = map(any)
  description = "List of our machines to be created in cloud"
}
