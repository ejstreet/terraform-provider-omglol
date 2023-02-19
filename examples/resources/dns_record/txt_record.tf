resource omglol_dns_record txt {
  type = "TXT"
  address = "example"
  name = "txt"
  data = "terraform=true"
  ttl = 300
}
