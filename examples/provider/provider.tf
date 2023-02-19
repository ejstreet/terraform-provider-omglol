terraform {
  required_providers {
    omglol = {
      source = "ejstreet/omglol"
    }
  }
}

provider omglol {
  api_key = "<omg.lol API key>"
  email = "<email address>"
}
