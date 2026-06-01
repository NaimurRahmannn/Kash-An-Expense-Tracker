package controllers

// HealthController handles health check requests.
type HealthController struct {
	BaseController
}

// Get returns the API health status.
// @Title HealthCheck
// @Summary Check server health
// @Description Returns a basic success response when the API server is running.
// @Success 200 {object} controllers.SuccessResponse "Server is running"
// @Failure 500 {object} controllers.ErrorResponse "Internal server error"
// @router /health [get]
func (c *HealthController) Get() {
	c.Success("Server is running")
}
