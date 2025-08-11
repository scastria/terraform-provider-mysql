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
}
