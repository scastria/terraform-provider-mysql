package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		DeleteContext: resourceUserRoleDelete,
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
				ForceNew: true,
			},
		},
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	user := d.Get("user").(string)
	role := d.Get("role").(string)
	query, _, err := c.Exec(ctx, "grant '%s' to '%s'", role, user)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId(fmt.Sprintf("%s:%s", user, role))
	return diags
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	tokens := strings.Split(d.Id(), ":")
	user := tokens[0]
	role := tokens[1]
	// Must check for existence of user first
	var rowUser string
	query, row := c.QueryRow(ctx, "select user from mysql.user where user = '%s' and host = '%%'", user)
	err := row.Scan(&rowUser)
	if err != nil {
		d.SetId("")
		if errors.Is(err, sql.ErrNoRows) {
			return diags
		}
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	query, rows, err := c.Query(ctx, "show grants for '%s'", user)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	var foundPerm bool
	for rows.Next() {
		var rowPerm string
		err = rows.Scan(&rowPerm)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if strings.Contains(strings.ToLower(rowPerm), "grant") && strings.Contains(rowPerm, user) && strings.Contains(rowPerm, role) {
			foundPerm = true
			break
		}
	}
	if !foundPerm {
		d.SetId("")
		return diags
	}
	d.Set("user", user)
	d.Set("role", role)
	return diags
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	tokens := strings.Split(d.Id(), ":")
	user := tokens[0]
	role := tokens[1]
	query, _, err := c.Exec(ctx, "revoke '%s' from '%s'", role, user)
	if err != nil {
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId("")
	return diags
}
