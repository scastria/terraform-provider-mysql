package mysql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MYSQL_HOST", nil),
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MYSQL_PORT", 3306),
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MYSQL_DATABASE", "information_schema"),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MYSQL_USERNAME", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("MYSQL_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mysql_user":              resourceUser(),
			"mysql_user_role":         resourceUserRole(),
			"mysql_user_default_role": resourceUserDefaultRole(),
			"mysql_role":              resourceRole(),
			"mysql_role_permission":   resourceRolePermission(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("host").(string)
	port := d.Get("port").(int)
	database := d.Get("database").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics
	c, err := client.NewClient(host, port, database, username, password)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
