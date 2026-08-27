# Data Source: mysql_role
Represents a role
## Example usage
```hcl
data "mysql_role" "example" {
  search_name = "admin"
}
```
## Argument Reference
* `search_name` - **(Optional, String)** The search string to apply to the name of the role. Uses contains.
* `name` - **(Optional, String)** The filter string to apply to the name of the role. Uses equality.
## Attribute Reference
* `id` - **(String)** Same as `name`
