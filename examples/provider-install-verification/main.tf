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

resource "omglol_account_settings" "this" {
  communication = "email_ok"
  date_format = "iso_8601"
}

output info {
  value = data.omglol_account_info.this
}

output settings {
  value = omglol_account_settings.this
}

