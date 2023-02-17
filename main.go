package main

import (
    "context"
    "terraform-provider-omglol/omglol"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name omglol

func main() {
    providerserver.Serve(context.Background(), omglol.New, providerserver.ServeOpts{
        Address: "registry.terraform.io/ejstreet/omglol",
    })
}
