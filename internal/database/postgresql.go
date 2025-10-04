package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"varaden/server/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func connectPostgresql(cfg config.DBConfig) (*sql.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.Timezone)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Database connected: Postgresql")

	return db, nil
}
