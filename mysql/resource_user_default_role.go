package mysql

import (
	"context"
	"database/sql"
	"errors"

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
	query, _, err := c.Exec(ctx, "set default role '%s' to '%s'", role, user)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId(user)
	return diags
}

func resourceUserDefaultRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	user := d.Id()
	var defaultRole string
	query, row := c.QueryRow(ctx, "select default_role_user from mysql.default_roles where user = '%s' and host = '%%' and default_role_host = '%%'", user)
	err := row.Scan(&defaultRole)
	if err != nil {
		d.SetId("")
		if errors.Is(err, sql.ErrNoRows) {
			return diags
		}
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.Set("user", user)
	d.Set("role", defaultRole)
	return diags
}

func resourceUserDefaultRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	user := d.Id()
	role := d.Get("role").(string)
	query, _, err := c.Exec(ctx, "set default role '%s' to '%s'", role, user)
	if err != nil {
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	return diags
}

func resourceUserDefaultRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	user := d.Id()
	query, _, err := c.Exec(ctx, "set default role NONE to '%s'", user)
	if err != nil {
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId("")
	return diags
}
