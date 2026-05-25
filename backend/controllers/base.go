package controllers

import beego "github.com/beego/beego/v2/server/web"

// BaseController provides shared API response helpers.
type BaseController struct {
	beego.Controller
}

type apiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a successful JSON response without data.
func (c *BaseController) Success(message string) {
	c.SuccessWithStatus(200, message)
}

// SuccessWithStatus sends a successful JSON response with the provided HTTP status code.
func (c *BaseController) SuccessWithStatus(statusCode int, message string) {
	c.Ctx.Output.SetStatus(statusCode)
	c.Data["json"] = apiResponse{
		Success: true,
		Message: message,
	}
	c.ServeJSON()
}

// SuccessWithData sends a successful JSON response with data.
func (c *BaseController) SuccessWithData(message string, data interface{}) {
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = apiResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.ServeJSON()
}

// ErrorResponse sends an error JSON response with the provided HTTP status code.
func (c *BaseController) ErrorResponse(statusCode int, message string) {
	c.Ctx.Output.SetStatus(statusCode)
	c.Data["json"] = apiResponse{
		Success: false,
		Message: message,
	}
	c.ServeJSON()
}
