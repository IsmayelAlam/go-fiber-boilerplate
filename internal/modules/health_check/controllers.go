package healthCheck

import (
	"github.com/gofiber/fiber/v2"
)

// Health check
//
//	@Summary		Application health check
//	@Description	Checks the status of critical dependencies (e.g., database, memory). Returns overall health status and individual service statuses.
//	@Tags			System
//	@Produce		json
//	@Success		200	{object}	HealthCheckResponse	"All services are healthy"
//	@Success		500	{object}	utils.CommonError	"One or more services are unhealthy"
//	@Router			/health-check [get]
func (s *healthCheckService) Check(c *fiber.Ctx) error {
	isHealthy := true
	var serviceList []HealthCheck

	// Check the database connection
	if err := s.dbCheck(); err != nil {
		isHealthy = false
		errMsg := err.Error()
		s.addServiceStatus(&serviceList, "Postgresql", false, &errMsg)
	} else {
		s.addServiceStatus(&serviceList, "Postgresql", true, nil)
	}

	if err := s.memoryHeapCheck(); err != nil {
		isHealthy = false
		errMsg := err.Error()
		s.addServiceStatus(&serviceList, "Memory", false, &errMsg)
	} else {
		s.addServiceStatus(&serviceList, "Memory", true, nil)
	}

	// Return the response based on health check result
	statusCode := fiber.StatusOK
	status := "success"

	if !isHealthy {
		statusCode = fiber.StatusInternalServerError
		status = "error"
	}

	return c.Status(statusCode).JSON(HealthCheckResponse{
		Status:    status,
		Message:   "Health check completed",
		Code:      statusCode,
		IsHealthy: isHealthy,
		Result:    serviceList,
	})
}
