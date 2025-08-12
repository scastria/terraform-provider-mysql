package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_plugin": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"auth_plugin_alias"},
			},
			"auth_plugin_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"auth_plugin"},
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Get("name").(string)
	auth := ""
	authPlugin, ok := d.GetOk("auth_plugin")
	if ok {
		authPluginAlias := d.Get("auth_plugin_alias").(string)
		auth = fmt.Sprintf("identified with %s as '%s'", authPlugin, authPluginAlias)
	}
	query, _, err := c.Exec(ctx, "create user '%s' %s", name, auth)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId(name)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Id()
	var rowUser, rowPlugin, rowAuth string
	query, row := c.QueryRow(ctx, "select user, plugin, authentication_string from mysql.user where user = '%s' and host = '%%'", name)
	err := row.Scan(&rowUser, &rowPlugin, &rowAuth)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	// Get default authentication plugin
	var rowVar, rowDefaultPlugin string
	query, row = c.QueryRow(ctx, "show variables like 'default_authentication_plugin'")
	err = row.Scan(&rowVar, &rowDefaultPlugin)
	if err != nil {
		d.SetId("")
		if err == sql.ErrNoRows {
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("name", rowUser)
	if rowPlugin != rowDefaultPlugin {
		d.Set("auth_plugin", rowPlugin)
		d.Set("auth_plugin_alias", rowAuth)
	}
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		query, _, err := c.Exec(ctx, "rename user '%s' to '%s'", oldName.(string), newName.(string))
		if err != nil {
			return diag.Errorf("Error executing query: %s, error: %v", query, err)
		}
		d.SetId(newName.(string))
	}
	name := d.Id()
	auth := ""
	authPlugin, ok := d.GetOk("auth_plugin")
	if ok {
		authPluginAlias := d.Get("auth_plugin_alias").(string)
		auth = fmt.Sprintf("identified with %s as '%s'", authPlugin, authPluginAlias)
	} else {
		// Get default auth plugin if not specified
		var rowVar, rowVal string
		query, row := c.QueryRow(ctx, "show variables like 'default_authentication_plugin'")
		err := row.Scan(&rowVar, &rowVal)
		if err != nil {
			return diag.Errorf("Error executing query: %s, error: %v", query, err)
		}
		auth = fmt.Sprintf("identified with %s", rowVal)
	}
	query, _, err := c.Exec(ctx, "alter user '%s' %s", name, auth)
	if err != nil {
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Id()
	query, _, err := c.Exec(ctx, "drop user '%s'", name)
	if err != nil {
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId("")
	return diags
}
