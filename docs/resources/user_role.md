# Resource: mysql_user_role
Represents a role assigned to a user
## Example usage
```hcl
resource "mysql_user" "User" {
  name = "MyUser"
}
resource "mysql_role" "Role" {
  name = "MyRole"
}
resource "mysql_user_role" "example" {
  user = mysql_user.User.id
  role = mysql_role.Role.id
}
```
## Argument Reference
* `user` - **(Required, ForceNew, String)** The name of the user.
* `role` - **(Required, ForceNew, String)** The name of the role.
## Attribute Reference
* `id` - **(String)** Same as `user`:`role`
## Import
User roles can be imported using a proper value of `id` as described above
