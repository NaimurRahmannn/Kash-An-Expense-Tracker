package routers

import (
	"backend/controllers"
	"backend/docs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
)

// @Title Expense Tracker API
// @Description REST API for personal expense tracking with CSV local storage and optional Postgres production storage.
// @APIVersion 1.0.0
// @Schemes http
// @SecurityDefinition XUserID apiKey X-User-ID header "Use user_id returned by login as the X-User-ID header for expense endpoints."
func init() {
	if swaggerEnabled() {
		docs.RegisterSwaggerRoutes()
	}

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://kash-an-expense-tracker.vercel.app",
		},
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

func swaggerEnabled() bool {
	enabled, err := beego.AppConfig.Bool("enable_swagger")
	return err == nil && enabled
}

func swaggerDocsNamespace() {
	_ = beego.NewNamespace("/api/v1",
		beego.NSInclude(
			&controllers.HealthController{},
			&controllers.AuthController{},
			&controllers.ExpenseController{},
		),
	)
}
