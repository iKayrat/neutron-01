package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
	log.Println("Env $PORT :", os.Getenv("PORT"))
	if os.Getenv("PORT") != "" {
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatal(err)
			log.Fatal("$PORT must be set")
		}
		log.Println("port : ", port)
		beego.BConfig.Listen.HTTPPort = port
		beego.BConfig.Listen.HTTPSPort = port
	}
	beego.Run()
}
