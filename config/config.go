package config

import (
	"flag"
	"fmt"
)

type DBConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
	Timezone string

	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  int
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type AllConfig struct {
	PortAddress string
	DB          DBConfig
	SMTP        SMTPConfig
}

func AppConfig() AllConfig {

	var cfg AllConfig

	// App config
	flag.IntVar(&Port, "port", 8080, "API server port")
	flag.Func("env", "Environment (development|staging|production)", func(s string) error {
		switch s {
		case "development":
			IsDevelopment = true
			return nil
		case "staging":
			IsStaging = true
			return nil
		case "production":
			IsProduction = true
			return nil
		default:
			return fmt.Errorf("invalid environment %q, must be one of: development, staging, production", s)
		}
	})

	// Database config
	flag.StringVar(&cfg.DB.Host, "db-host", "localhost", "PostgreSQL host name")
	flag.StringVar(&cfg.DB.Port, "db-port", "5432", "PostgreSQL port number")
	flag.StringVar(&cfg.DB.User, "db-user", "postgres", "PostgreSQL user")
	flag.StringVar(&cfg.DB.Password, "db-password", "postgres", "PostgreSQL password")
	flag.StringVar(&cfg.DB.DBName, "db-dbname", "postgres", "PostgreSQL database name")
	flag.StringVar(&cfg.DB.SSLMode, "db-sslmode", "disable", "PostgreSQL SSL mode")
	flag.StringVar(&cfg.DB.Timezone, "db-timezone", "UTC", "PostgreSQL DSN")

	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.IntVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15, "PostgreSQL max connection idle time in minutes")

	// SMTP config
	flag.StringVar(&cfg.SMTP.Host, "smtp-host", "smtp.example.com", "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", "user@example.com", "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", "password", "SMTP password")
	flag.StringVar(&cfg.SMTP.From, "smtp-from", "noreply@example.com", "SMTP from address")

	// set constance
	flag.StringVar(&FrontEndURL, "frontend-url", "http://localhost:3000", "Front end URL")

	// JWT-Config
	flag.StringVar(&JWTConfig.Issuer, "jwt-issuer", "myapp.example.com", "JWT Issuer (typically your service domain)")
	flag.StringVar(&JWTConfig.Secret, "jwt-secret", "1234567890", "JWT Secret (should be a strong, random secret; leave empty to require via env/config)")
	flag.IntVar(&JWTConfig.TokenExpiry, "jwt-token-expiry", 6, "JWT Access Token Expiry in hours")
	flag.IntVar(&JWTConfig.RefreshExpiry, "jwt-refresh-expiry", 30, "JWT Refresh Token Expiry in days")
	flag.StringVar(&JWTConfig.RefreshCookieDomain, "jwt-cookie-domain", "localhost", "JWT Cookie Domain (include leading dot for subdomain sharing)")
	flag.StringVar(&JWTConfig.RefreshCookiePath, "jwt-cookie-path", "/", "JWT Cookie Path")
	flag.StringVar(&JWTConfig.RefreshCookieName, "jwt-cookie-name", "__r_token", "JWT Cookie Name")

	flag.Parse()

	JWTConfig.Audience = FrontEndURL

	cfg.PortAddress = fmt.Sprintf(":%d", Port)

	return cfg
}
