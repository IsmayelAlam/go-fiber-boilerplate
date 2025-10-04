package config

import (
	"varaden/server/internal/utils"
)

var (
	FrontEndURL   = ""
	Port          int
	IsDevelopment = false
	IsStaging     = false
	IsProduction  = false
	JWTConfig     = &utils.JWTConfig{
		Issuer:              "your-issuer",
		Audience:            "your-audience",
		Secret:              "your-secret",
		TokenExpiry:         6,
		RefreshExpiry:       30,
		RefreshCookieDomain: "your-cookie-domain",
		RefreshCookiePath:   "your-cookie-path",
		RefreshCookieName:   "your-cookie-name",
	}
)
