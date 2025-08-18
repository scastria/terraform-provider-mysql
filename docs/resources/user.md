# Resource: mysql_user
Represents a user
## Example usage
```hcl
resource "mysql_user" "example" {
  name = "MyUser"
}
```
## Argument Reference
* `name` - **(Required, String)** The name of the user.
* `auth_plugin` - **(Optional, String)** The plugin to use for authentication.
* `auth_plugin_alias` - **(Optional, String)** The string used by the auth plugin. Must be specified if `auth_plugin` is specified.
* `email` - **(Optional, String)** The email of the user stored in metadata attributes.
## Attribute Reference
* `id` - **(String)** Same as `name`.
## Import
Users can be imported using a proper value of `id` as described above
