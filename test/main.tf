terraform {
  required_providers {
    mysql = {
      source = "github.com/scastria/mysql"
    }
  }
}

provider "mysql" {
}

resource "mysql_user" "User" {
  name = "TestUser"
  auth_plugin = "AWSAuthenticationPlugin"
  auth_plugin_alias = "RDS"
  email = "good@bad.com"
}
resource "mysql_role" "Role" {
  name = "TestRole"
}
# resource "mysql_user_role" "UserRole" {
#   user = mysql_user.User.name
#   role = mysql_role.Role.name
# }
# resource "mysql_user_default_role" "UserDefaultRole" {
#   user = mysql_user.User.name
#   role = mysql_role.Role.name
# }
resource "mysql_role_permission" "RolePermission" {
  role = mysql_role.Role.name
  privilege = "CREATE"
}
