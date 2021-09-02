package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"neutron0.1/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/todo", &controllers.MainController{}, "get,post:CreateToDo")
	beego.Router("/api/register", &controllers.UserController{}, "get,post:Register")
	beego.Router("/api/login", &controllers.UserController{}, "get,post:Login")
	beego.Router("/api/logout", &controllers.UserController{}, "get,post:Logout")
	beego.Router("/token/refresh", &controllers.UserController{}, "get,post:Refresh")

	beego.Router("/login", &controllers.UserController{}, "get,post:Loginsession")
	beego.Router("/auth", &controllers.UserController{}, "get,post:Auth")

	// beego.Get("/", func(ctx *context.Context) {
	// 	_ = ctx.Output.Body([]byte("This is a Beego + JWT API - Creator: Mehran Abghari (mehran.ab80@gmail.com)"))
	// })

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)

}
