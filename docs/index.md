# MySQL Provider
The MySQL provider is used to manage various administrative resources.  The provider
needs to be configured with the proper credentials before it can be used.

This provider does NOT cover 100% of the MySQL features.  If there is something missing
that you would like to be added, please submit an Issue in corresponding GitHub repo.
## Example Usage
```hcl
terraform {
  required_providers {
    konnect = {
      source  = "scastria/mysql"
      version = "~> 0.1.0"
    }
  }
}

# Configure the MySQL Provider
provider "mysql" {
  host = "myserver.example.com"
  username = "XXXXX"
  password = "YYYYY"
}
```
## Argument Reference
* `host` - **(Required, String)** The hostname of the mysql server. Can be specified via env variable `MYSQL_HOST`.
* `port` - **(Optional, Integer)** The port of the mysql server. Can be specified via env variable `MYSQL_PORT`. Default: `3306`
* `database` - **(Optional, String)** The default database/schema to connect to. Can be specified via env variable `MYSQL_DATABASE`. Default: `information_schema`.
* `username` - **(Required, String)** Username to connect to server as. Can be specified via env variable `MYSQL_USERNAME`.
* `password` - **(Required, String)** Password to connect to server with. Can be specified via env variable `MYSQL_PASSWORD`.
