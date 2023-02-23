terraform {
  required_providers {
    omglol = {
      source = "ejstreet/omglol"
      version = "<version>"
    }
  }
}

provider omglol {
  api_key = "<omg.lol API key>"
  user_email = "<email address>"
}
