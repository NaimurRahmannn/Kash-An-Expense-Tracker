package routers

import (
	"backend/controllers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Accept",
			"Content-Type",
			"Authorization",
			"X-User-ID",
		},
	}))

	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/health", &controllers.HealthController{}, "get:Get"),
		beego.NSNamespace("/auth",
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),
		beego.NSRouter("/expenses", &controllers.ExpenseController{}, "get:List;post:Create"),
		beego.NSRouter("/expenses/summary", &controllers.ExpenseController{}, "get:Summary"),
		beego.NSRouter("/expenses/:id", &controllers.ExpenseController{}, "get:GetOne;put:Update;delete:Delete"),
	)

	beego.AddNamespace(ns)
}
