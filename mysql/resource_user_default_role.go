package mysql

import (
	"context"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func resourceUserDefaultRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserDefaultRoleCreate,
		ReadContext:   resourceUserDefaultRoleRead,
		UpdateContext: resourceUserDefaultRoleUpdate,
		DeleteContext: resourceUserDefaultRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceUserDefaultRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	user := d.Get("user").(string)
	role := d.Get("role").(string)
	db, err := c.DbConnection()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	_, err = c.Exec(ctx, db, "set default role '%s' to '%s'", role, user)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s", user, role))
	return diags
}

func resourceUserDefaultRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	tokens := strings.Split(d.Id(), ":")
	user := tokens[0]
	role := tokens[1]
	db, err := c.DbConnection()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	var count int
	err = c.QueryRow(ctx, db, "select count(*) from mysql.default_roles where user = '%s' and host = '%%' and default_role_user = '%s' and default_role_host = '%%'", user, role).Scan(&count)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	if count == 0 {
		d.SetId("")
		return diags
	}
	d.Set("user", user)
	d.Set("role", role)
	return diags
}

func resourceUserDefaultRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	user := d.Get("user").(string)
	role := d.Get("role").(string)
	db, err := c.DbConnection()
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	defer db.Close()
	_, err = c.Exec(ctx, db, "set default role '%s' to '%s'", role, user)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s:%s", user, role))
	return diags
}

func resourceUserDefaultRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	tokens := strings.Split(d.Id(), ":")
	user := tokens[0]
	db, err := c.DbConnection()
	if err != nil {
		return diag.FromErr(err)
	}
	defer db.Close()
	_, err = c.Exec(ctx, db, "set default role NONE to '%s'", user)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
