// Package docs serves the generated Swagger specification and Swagger UI.
package docs

import (
	_ "embed"

	beecontext "github.com/beego/beego/v2/server/web/context"

	beego "github.com/beego/beego/v2/server/web"
)

//go:embed swagger.json
var swaggerJSON []byte

// RegisterSwaggerRoutes registers local Swagger documentation routes.
func RegisterSwaggerRoutes() {
	beego.Get("/swagger/", serveSwaggerUI)
	beego.Get("/swagger/index.html", serveSwaggerUI)
	beego.Get("/swagger/doc.json", serveSwaggerJSON)
}

func serveSwaggerUI(ctx *beecontext.Context) {
	ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	_ = ctx.Output.Body([]byte(swaggerHTML))
}

func serveSwaggerJSON(ctx *beecontext.Context) {
	ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	_ = ctx.Output.Body(swaggerJSON)
}

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Expense Tracker API Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      window.ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui"
      });
    };
  </script>
</body>
</html>`
