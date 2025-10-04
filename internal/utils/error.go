package utils

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Common struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ErrorDetails struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

var customMessages = map[string]string{
	"required": "Field %s must be filled",
	"email":    "Invalid email address for field %s",
	"min":      "Field %s must have a minimum length of %s characters",
	"max":      "Field %s must have a maximum length of %s characters",
	"len":      "Field %s must be exactly %s characters long",
	"number":   "Field %s must be a number",
	"positive": "Field %s must be a positive number",
	"alphanum": "Field %s must contain only alphanumeric characters",
	"oneof":    "Invalid value for field %s",
	"password": "Field %s must contain at least 1 letter and 1 number",
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if errorsMap := CustomErrorMessages(err); len(errorsMap) > 0 {
		return Error(c, fiber.StatusBadRequest, "Bad Request test", errorsMap)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return Error(c, fiberErr.Code, fiberErr.Message, nil)
	}

	return Error(c, fiber.StatusInternalServerError, "Internal Server Error", nil)
}

func NotFoundHandler(c *fiber.Ctx) error {
	return Error(c, fiber.StatusNotFound, "Endpoint Not Found", nil)
}

func DuplicateEntryError(err error, field string) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("%s already in use", field))
	}
	return err
}

func Error(c *fiber.Ctx, statusCode int, message string, details any) error {
	var errRes error
	if details != nil {
		errRes = c.Status(statusCode).JSON(ErrorDetails{
			Code:    statusCode,
			Status:  "error",
			Message: message,
			Errors:  details,
		})
	} else {
		errRes = c.Status(statusCode).JSON(Common{
			Code:    statusCode,
			Status:  "error",
			Message: message,
		})
	}

	if errRes != nil {
		log.Errorf("Failed to send error response : %+v", errRes)
	}

	return errRes
}

func CustomErrorMessages(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return generateErrorMessages(validationErrors)
	}
	return nil
}

func generateErrorMessages(validationErrors validator.ValidationErrors) map[string]string {
	errorsMap := make(map[string]string)
	for _, err := range validationErrors {
		fieldName := err.StructNamespace()
		tag := err.Tag()

		customMessage := customMessages[tag]
		if customMessage != "" {
			errorsMap[fieldName] = formatErrorMessage(customMessage, err, tag)
		} else {
			errorsMap[fieldName] = defaultErrorMessage(err)
		}
	}
	return errorsMap
}
func formatErrorMessage(customMessage string, err validator.FieldError, tag string) string {
	if tag == "min" || tag == "max" || tag == "len" {
		return fmt.Sprintf(customMessage, err.Field(), err.Param())
	}
	return fmt.Sprintf(customMessage, err.Field())
}

func defaultErrorMessage(err validator.FieldError) string {
	return fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag())
}
