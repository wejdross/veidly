### Hetzner connector

hcloud_token      = "ywIGUGk8YMJMUnsH9pZvLkLEtZxaoO543dPEfyTtR13gjo7VZRTU4Eiu73e8SXAZ"
hetzner_dns_token = "P8zvC3g4z0XzT3B10TiizQRJkMDkV00u"
### 

vms = {
  "git" : {
    server_type      = "cpx31"
    data_volume_size = 50
  },
  back : {
    data_volume_size = 100
  },
  rc : {
    server_type      = "cpx21"
    data_volume_size = 20
  },
  cass1-dev: {
    data_volume_size = 11,
  },
  cass2-dev: {
    data_volume_size = 11,
  },
  cass1-prod: {
    data_volume_size = 11,
  },
  cass2-prod: {
    data_volume_size = 11,
  },
  "app-prod": {
    server_type      = "cpx31"
    data_volume_size = 50
  },
  "app-dev": {
    server_type = "cpx31"
    data_volume_size = 11
  }
}
