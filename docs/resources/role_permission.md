# Resource: mysql_role_permission
Represents a permission of a role
## Example usage
```hcl
resource "mysql_role" "Role" {
  name = "MyRole"
}
resource "mysql_role_permission" "example" {
  role    = mysql_role.Role.id
  privilege  = "CREATE"
}
```
## Argument Reference
* `role` - **(Required, ForceNew, String)** The name of the role.
* `privilege` - **(Required, ForceNew, String)** The privilege to grant. Must be a valid privilege for MySQL.
* `level` - **(Optional, ForceNew, String)** At what level to grant the `privilege`. Allowed values: `global`, `database`, `table`, `function`, `procedure`. Default: `global`.
* `target` - **(Optional, ForceNew, String)** The target of the `privilege`. Must be specified when `level` is NOT `global`. 
## Attribute Reference
* `id` - **(String)** Same as `role`:`privilege`:`level`:`target`. Use empty string for parts of the id that do not apply.
## Import
Role permissions can be imported using a proper value of `id` as described above
