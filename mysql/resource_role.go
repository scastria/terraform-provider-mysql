package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRoleCreate,
		ReadContext:   resourceRoleRead,
		UpdateContext: resourceRoleUpdate,
		DeleteContext: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Get("name").(string)
	db, err := c.DbConnection()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	_, err = c.Exec(ctx, db, "create role '%s'", name)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(name)
	return diags
}

func resourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Id()
	db, err := c.DbConnection()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	var rowUser string
	err = c.QueryRow(ctx, db, "select User from mysql.user where User = '%s' and Host = '%%'", name).Scan(&rowUser)
	if err != nil {
		d.SetId("")
		if err == sql.ErrNoRows {
			return diags
		}
		return diag.FromErr(err)
	}
	d.Set("name", rowUser)
	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	db, err := c.DbConnection()
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		_, err = c.Exec(ctx, db, "rename user '%s' to '%s'", oldName.(string), newName.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(newName.(string))
	}
	return diags
}

func resourceRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	name := d.Id()
	db, err := c.DbConnection()
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	_, err = c.Exec(ctx, db, "drop role '%s'", name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
