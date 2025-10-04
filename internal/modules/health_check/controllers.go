package healthCheck

import (
	"github.com/gofiber/fiber/v2"
)

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
