package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/scastria/terraform-provider-mysql/mysql/client"
)

func resourceRolePermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRolePermissionCreate,
		ReadContext:   resourceRolePermissionRead,
		DeleteContext: resourceRolePermissionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"privilege": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"target"},
				Default:      "global",
				ValidateFunc: validation.StringInSlice([]string{"global", "database", "table", "function", "procedure"}, false),
			},
			"target": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"level"},
			},
		},
	}
}

func translateTarget(level string, target string) string {
	if level == "global" {
		return "*.*"
	}
	if level == "database" {
		return fmt.Sprintf("%s.*", target)
	}
	if level == "table" {
		return target
	}
	if level == "function" || level == "procedure" {
		return fmt.Sprintf("%s %s", level, target)
	}
	return ""
}

func resourceRolePermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	role := d.Get("role").(string)
	privilege := d.Get("privilege").(string)
	level := d.Get("level").(string)
	targetRaw, ok := d.GetOk("target")
	var target string
	if ok {
		target = targetRaw.(string)
	} else {
		target = ""
	}
	on := translateTarget(level, target)
	query, _, err := c.Exec(ctx, "grant %s on %s to '%s'", privilege, on, role)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId(fmt.Sprintf("%s:%s:%s:%s", role, privilege, level, target))
	return diags
}

func resourceRolePermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	tokens := strings.Split(d.Id(), ":")
	role := tokens[0]
	privilege := tokens[1]
	level := tokens[2]
	target := tokens[3]
	on := translateTarget(level, target)
	query, rows, err := c.Query(ctx, "show grants for '%s'", role)
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	var foundPerm bool
	re := regexp.MustCompile(`GRANT\s+(.+)\s+ON\s+(.+)\s+TO\s+.*`)
	onLower := strings.ToLower(on)
	privilegeLower := strings.ToLower(privilege)
	for rows.Next() {
		var rowPerm string
		err = rows.Scan(&rowPerm)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		matches := re.FindStringSubmatch(rowPerm)
		if matches == nil {
			d.SetId("")
			return diag.Errorf("Unable to parse grant statement: %s", rowPerm)
		}
		// First check target since grant option priv is special
		if strings.Contains(strings.ReplaceAll(strings.ToLower(matches[2]), "`", ""), onLower) {
			if ((privilegeLower == "grant option") && strings.Contains(strings.ToLower(rowPerm), "with grant option")) || (strings.Contains(strings.ToLower(matches[1]), privilegeLower)) {
				foundPerm = true
				break
			}
		}
	}
	if !foundPerm {
		d.SetId("")
		return diags
	}
	d.Set("role", role)
	d.Set("privilege", privilege)
	d.Set("level", level)
	if target != "" {
		d.Set("target", target)
	}
	return diags
}

func resourceRolePermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	tokens := strings.Split(d.Id(), ":")
	role := tokens[0]
	privilege := tokens[1]
	level := tokens[2]
	target := tokens[3]
	on := translateTarget(level, target)
	query, _, err := c.Exec(ctx, "revoke %s on %s from '%s'", privilege, on, role)
	if err != nil {
		return diag.Errorf("Error executing query: %s, error: %v", query, err)
	}
	d.SetId("")
	return diags
}
