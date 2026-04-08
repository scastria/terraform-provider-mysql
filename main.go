package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scastria/terraform-provider-mysql/mysql"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mysql.Provider,
	})
	//c, err := client.NewClient(os.Getenv("HOST"), 3306, "information_schema", os.Getenv("USERNAME"), os.Getenv("PASSWORD"), 2, 2)
	//if err != nil {
	//	panic(err)
	//}
	//printStats(*c)
	//var ctx = context.Background()
	//var count int
	//_, row := c.QueryRow(ctx, "select count(*) from mysql.user where user = '%s' and host = '%%'", "TestRole1")
	//printStats(*c)
	//err = row.Scan(&count)
	//printStats(*c)
	//if err != nil {
	//	panic(err)
	//}
	//if count == 0 {
	//	println("Role does not exist")
	//}
	//println("Role count:", count)
	//_, row = c.QueryRow(ctx, "select count(*) from mysql.user where user = '%s' and host = '%%'", "TestRole1")
	//printStats(*c)
	//err = row.Scan(&count)
	//printStats(*c)
	//if err != nil {
	//	panic(err)
	//}
	//if count == 0 {
	//	println("Role does not exist")
	//}
	//println("Role count:", count)
	//_, row = c.QueryRow(ctx, "select count(*) from mysql.user where user = '%s' and host = '%%'", "TestRole1")
	//printStats(*c)
	//err = row.Scan(&count)
	//printStats(*c)
	//if err != nil {
	//	panic(err)
	//}
	//if count == 0 {
	//	println("Role does not exist")
	//}
	//println("Role count:", count)
	//printStats(*c)
}

//func printStats(c client.Client) {
//	var stats = c.Conn.Stats()
//	println("Stats:", "InUse:", stats.InUse, "Idle:", stats.Idle, "Open:", stats.OpenConnections)
//}
