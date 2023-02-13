terraform {
  required_providers {
    omglol = {
      source = "ejstreet/omglol"
    }
  }
}

provider "omglol" {}

data "omglol_account_info" "this" {}

output info {
  value = formatdate("EEEE, DD-MMM-YY hh:mm:ss ZZZ", data.omglol_account_info.this.created)
}