package mysql

import (
	"context"

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
	_, err := c.Exec(ctx, "create role '%s'", name)
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
	var count int
	err := c.QueryRow(ctx, "select count(*) from mysql.user where user = '%s' and host = '%%'", name).Scan(&count)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if count == 0 {
		d.SetId("")
		return diags
	}
	d.Set("name", name)
	return diags
}

func resourceRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if d.HasChange("name") {
		oldName, newName := d.GetChange("name")
		_, err := c.Exec(ctx, "rename user '%s' to '%s'", oldName.(string), newName.(string))
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
	_, err := c.Exec(ctx, "drop role '%s'", name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
