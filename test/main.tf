terraform {
  required_providers {
    mysql = {
      source = "github.com/scastria/mysql"
    }
  }
}

provider "mysql" {
  max_open_connections = 2
  max_idle_connections = 2
}

locals {
  roles = [
    "TestRole1",
    "TestRole2",
    "TestRole3",
    "TestRole4",
    "TestRole5",
    "TestRole6",
    "TestRole7",
    "TestRole8",
    "TestRole9",
    "TestRole10"
  ]
}
# resource "mysql_user" "User" {
#   name = "TestUser"
#   auth_plugin = "AWSAuthenticationPlugin"
#   auth_plugin_alias = "RDS"
#   email = "good@bad.com"
# }
resource "mysql_role" "Role" {
  for_each = toset(local.roles)
  name = each.key
}
# resource "mysql_role" "Role2" {
#   name = "TestRole2"
# }
# resource "mysql_user_role" "UserRole" {
#   user = mysql_user.User.name
#   role = mysql_role.Role.name
# }
# resource "mysql_user_role" "UserRole2" {
#   user = mysql_user.User.name
#   role = mysql_role.Role2.name
# }
# resource "mysql_user_default_role" "UserDefaultRole" {
#   user = mysql_user.User.name
#   role = mysql_role.Role.name
# }
# resource "mysql_role_permission" "RolePermission" {
#   role = mysql_role.Role.name
#   privilege = "CREATE"
# }
# resource "mysql_role_permission" "RolePermission2" {
#   role = mysql_role.Role.name
#   privilege = "GRANT OPTION"
# }
