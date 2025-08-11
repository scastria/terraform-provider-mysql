terraform {
  required_providers {
    mysql = {
      source = "github.com/scastria/mysql"
    }
  }
}

provider "mysql" {
  host = ""
  username = ""
  password = ""
}

resource "mysql_user" "User" {
  name = "TestUser2"
  auth_plugin = "AWSAuthenticationPlugin"
  auth_plugin_alias = "RDS"
}
