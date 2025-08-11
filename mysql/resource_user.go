package mysql

import (
	"context"
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
	db, err := c.DbConnection(ctx)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	auth := ""
	authPlugin, ok := d.GetOk("auth_plugin")
	if ok {
		authPluginAlias := d.Get("auth_plugin_alias").(string)
		auth = fmt.Sprintf("identified with %s as '%s'", authPlugin, authPluginAlias)
	}
	_, err = db.Exec(fmt.Sprintf("create user '%s' %s", name, auth))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(name)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Id()
	db, err := c.DbConnection(ctx)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	var rowUser, rowPlugin, rowAuth string
	err = db.QueryRow(fmt.Sprintf("select User, plugin, authentication_string from mysql.user where User = '%s' and Host = '%%'", name)).Scan(&rowUser, &rowPlugin, &rowAuth)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	// Get default authentication plugin
	var rowVar, rowDefaultPlugin string
	err = db.QueryRow("show variables like 'default_authentication_plugin'").Scan(&rowVar, &rowDefaultPlugin)
	if err != nil {
		d.SetId("")
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
	db, err := c.DbConnection(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		_, err = db.Exec(fmt.Sprintf("rename user '%s' to '%s'", oldName.(string), newName.(string)))
		if err != nil {
			return diag.FromErr(err)
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
		err = db.QueryRow("show variables like 'default_authentication_plugin'").Scan(&rowVar, &rowVal)
		if err != nil {
			return diag.FromErr(err)
		}
		auth = fmt.Sprintf("identified with %s", rowVal)
	}
	_, err = db.Exec(fmt.Sprintf("alter user '%s' %s", name, auth))
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Id()
	db, err := c.DbConnection(ctx)
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	_, err = db.Exec(fmt.Sprintf("drop user '%s'", name))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
