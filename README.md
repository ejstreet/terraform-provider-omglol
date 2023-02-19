# terraform-provider-omglol
Terraform omg.lol provider

## Getting Started
To use this provider, start by adding the provider definition:
```terraform
terraform {
  required_providers {
    omglol = {
      source = "ejstreet/omglol"
    }
  }
}

provider omglol {
  api_key = "<omg.lol API key>" // Alternatively set the OMGLOL_API_KEY env variable
  email = "<email address>" // Alternatively set the OMGLOL_USER_EMAIL env variable
}
```

You can quickly check your configuration by retrieving your account settings using the account settings data source:
```terraform
data omglol_account_info this {}
```

For further detail, see the provider documentation on the terraform registry.

## Local Development
Start by cloning the repo to your Go path. If forking the repo, update the path accordingly:
```bash
DEST=$GOROOT/src/github.com/ejstreet/terraform-provider-omglol/
mkdir -p $DEST
git clone git@github.com:ejstreet/terraform-provider-omglol.git $DEST
```

Then update your `~/.terraformrc` file with the following content, where `/home/user/go` is your `$GOROOT`:
```hcl
provider_installation {

  dev_overrides {
      "ejstreet/omglol" = "/home/user/go/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```
When in the root of the repo, run 
```bash
go install .
```
To install the project. Run this any time you make changes.
