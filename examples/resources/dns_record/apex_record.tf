resource omglol_dns_record apex {
  type = "TXT"
  address = "example"
  name = "@"
  data = "terraform=true"
  ttl = 300
}
