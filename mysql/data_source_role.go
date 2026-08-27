package mysql

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func dataSourceRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRoleRead,
		Schema: map[string]*schema.Schema{
			"search_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"search_name"},
			},
		},
	}
}

func dataSourceRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	var whereClause string
	searchName, ok := d.GetOk("search_name")
	if ok {
		whereClause = fmt.Sprintf("user like '%%%s%%'", searchName.(string))
	}
	name, ok := d.GetOk("name")
	if ok {
		whereClause = fmt.Sprintf("user = '%s'", name.(string))
	}
	query, rows, err := c.Query(ctx, "select user from mysql.user where %s and host = '%%'", whereClause)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	foundRole := false
	var foundName string
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&foundName)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		if foundRole {
			d.SetId("")
			return diag.FromErr(fmt.Errorf("Filter criteria does not result in a single role"))
		}
		foundRole = true
	}
	if !foundRole {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No role exists with that filter criteria"))
	}
	d.Set("name", foundName)
	d.SetId(foundName)
	return diags
}
