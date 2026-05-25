package controllers

// HealthController handles health check requests.
type HealthController struct {
	BaseController
}

// Get returns the API health status.
func (c *HealthController) Get() {
	c.Success("Server is running")
}
