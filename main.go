package main

import (
    "context"
    "terraform-provider-omglol/omglol"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
    providerserver.Serve(context.Background(), omglol.New, providerserver.ServeOpts{
        Address: "registry.terraform.io/ejstreet/omglol",
    })
}
