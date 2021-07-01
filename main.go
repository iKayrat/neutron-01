package main

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"

	_ "github.com/lib/pq"
	_ "neutron0.1/routers"
)

func init() {
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		fmt.Println("An error: ", err)
	}
}

func main() {
	orm.Debug = true

	beego.Run()
}
