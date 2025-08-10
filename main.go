package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-mysql/mysql"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mysql.Provider,
	})
}
