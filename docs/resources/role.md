# Resource: mysql_role
Represents a role
## Example usage
```hcl
resource "mysql_role" "example" {
  name = "MyRole"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the role.
## Attribute Reference
* `id` - **(String)** Same as `name`.
## Import
Roles can be imported using a proper value of `id` as described above
