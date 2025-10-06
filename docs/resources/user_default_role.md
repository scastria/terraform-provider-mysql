# Resource: mysql_user_default_role
Represents the default role to set as active when a user connects
## Example usage
```hcl
resource "mysql_user" "User" {
  name = "MyUser"
}
resource "mysql_role" "Role" {
  name = "MyRole"
}
resource "mysql_user_default_role" "example" {
  user = mysql_user.User.id
  role = mysql_role.Role.id
}
```
## Argument Reference
* `user` - **(Required, ForceNew, String)** The name of the user.
* `role` - **(Required, String)** The name of the default role. Must be a role that is assigned to the user.
## Attribute Reference
* `id` - **(String)** Same as `user`
## Import
User default roles can be imported using a proper value of `id` as described above
