package controllers

import (
	"encoding/json"
	"fmt"

	beego "github.com/beego/beego/v2/server/web"
	"neutron0.1/token"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

type ToDo struct {
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
}

func (c *MainController) CreateToDo() {
	var td *ToDo
	if err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &td); err != nil {
		c.Data["json"] = "invalid json"
		c.ServeJSON()
		c.StopRun()
		return
	}

	tokenAuth, err := newjwt.ExtractTokenMetadata(c.Ctx.Request)
	fmt.Println("*extract token tokenAuth:", tokenAuth)
	if err != nil {
		fmt.Println("json err:", err)
		c.Data["json"] = "Unauthorized extractmeta"
		c.ServeJSON()
		c.StopRun()
	}

	userid, err := token.FetchAuth(tokenAuth)
	if err != nil {
		c.Data["json"] = "unauthorized fetch"
		c.ServeJSON()
		c.StopRun()
		return
	}
	td.UserID = userid

	c.Data["json"] = td
	c.ServeJSON()
	c.StopRun()
}
