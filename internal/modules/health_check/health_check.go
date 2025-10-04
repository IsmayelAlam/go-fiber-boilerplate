package healthCheck

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type HealthCheck struct {
	Name    string  `json:"name"`
	Status  string  `json:"status"`
	IsUp    bool    `json:"is_up"`
	Message *string `json:"message,omitempty"`
}

type HealthCheckResponse struct {
	Code      int           `json:"code"`
	Status    string        `json:"status"`
	Message   string        `json:"message"`
	IsHealthy bool          `json:"is_healthy"`
	Result    []HealthCheck `json:"result"`
}

type healthCheckService struct {
	db    *sql.DB
	route fiber.Router
}

func RegisterHealthCheckModule(route fiber.Router, db *sql.DB) *healthCheckService {
	return &healthCheckService{
		db:    db,
		route: route,
	}
}
