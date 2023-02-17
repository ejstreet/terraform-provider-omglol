terraform {
  required_providers {
    omglol = {
      source = "ejstreet/omglol"
    }
  }
}

provider "omglol" {
  api_host = "https://api.omg.lol"
}

data "omglol_account_info" "this" {}

output info {
  value = data.omglol_account_info.this
}

resource "omglol_account_settings" "this" {
  communication = "email_ok"
  date_format = "iso_8601"
}

output settings {
  value = omglol_account_settings.this
}


resource "omglol_dns_record" "test" {
  type = "TXT"
  address = "terraform"
  name = "deployed"
  data = "terraform=true"
  ttl = 300
}

output "record" {
  value = omglol_dns_record.test
}

resource omglol_dns_record mx {
  type = "MX"
  address = "terraform"
  name = "mail"
  data = "mx_data"
  priority = 20
  ttl = 60
}

output mx {
  value = omglol_dns_record.mx
}

