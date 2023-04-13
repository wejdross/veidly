resource "hcloud_ssh_key" "terraform_key" {
  name       = "terraform-ssh"
  public_key = file("../../ssh/id.pub")
}

resource "hcloud_server" "veidly_infra_vms" {
  for_each    = var.vms
  name        = each.key
  image       = var.veidly_infra_machines.image
  server_type = lookup(each.value, "server_type", var.veidly_infra_machines.server_type)
  location    = "fsn1"
  keep_disk   = true
  labels = {
    terraform_managed = true,
  }
  ssh_keys = ["terraform-ssh"]
  depends_on = [
    hcloud_ssh_key.terraform_key
  ]
}
# volumeny na dane
resource "hcloud_volume" "veidly_infra_machines_volumes" {
  for_each  = hcloud_server.veidly_infra_vms
  name      = "${each.value.name}-data"
  server_id = each.value.id
  size      = lookup(var.vms[each.value.name], "data_volume_size", var.veidly_infra_machines.data_volume_size)
  depends_on = [
    hcloud_server.veidly_infra_vms
  ]
}
resource "hetznerdns_record" "veidly_record" {
  for_each = hcloud_server.veidly_infra_vms
  zone_id  = "7xYfmspYeirtFsu8m7dV3X"
  name     = "${each.value.name}.infra"
  value    = hcloud_server.veidly_infra_vms[each.key].ipv4_address
  type     = "A"
  depends_on = [
    hcloud_server.veidly_infra_vms
  ]
}

resource "hetznerdns_record" "veidly_prod" {
  zone_id  = "7xYfmspYeirtFsu8m7dV3X"
  name     = "@"
  value    = hcloud_server.veidly_infra_vms["app-prod"].ipv4_address
  type     = "A"
  depends_on = [
    hcloud_server.veidly_infra_vms
  ]
}

resource "hetznerdns_record" "veidly_dev" {
  zone_id  = "7xYfmspYeirtFsu8m7dV3X"
  name     = "dev"
  value    = hcloud_server.veidly_infra_vms["app-dev"].ipv4_address
  type     = "A"
  depends_on = [
    hcloud_server.veidly_infra_vms
  ]
}