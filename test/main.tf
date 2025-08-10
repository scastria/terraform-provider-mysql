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
