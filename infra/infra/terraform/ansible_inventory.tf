resource "local_file" "ansible_inventory" {
  content = templatefile("../ansible/ansible_inventory_template.j2",
    {
      prod_cluster = [
        for srv in hetznerdns_record.veidly_record : 
          length(regexall("^cass.*", srv.name)) > 0 ?
            "" :
            "${srv.name}.veidly.com ansible_host=${srv.value} terra_fqdn=${srv.name}.veidly.com"
        ],
        cass_dev_cluster = [
        for srv in hetznerdns_record.veidly_record : 
          length(regexall("^cass..dev*", srv.name)) > 0 ?
            "${srv.name}.veidly.com ansible_host=${srv.value} terra_fqdn=${srv.name}.veidly.com rack=1 dc=fsn1" : ""
        ]
        cass_prod_cluster = [
        for srv in hetznerdns_record.veidly_record : 
          length(regexall("^cass..prod*", srv.name)) > 0 ?
            "${srv.name}.veidly.com ansible_host=${srv.value} terra_fqdn=${srv.name}.veidly.com rack=1 dc=fsn1" : ""
        ]
    }
  )
  filename        = "../ansible/hosts.cfg"
  file_permission = "0600"
  depends_on = [
    hetznerdns_record.veidly_record
  ]
}
