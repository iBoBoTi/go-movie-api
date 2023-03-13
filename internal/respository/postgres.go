package respository

import (
	"context"
	"fmt"
	"github.com/iBoBoTi/go-movie-api/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// DB connection
func ConnectPostgres(c *config.Config) (*pgxpool.Pool, error) {
	log.Println("Connecting to Postgresql DB pool")
	dns := fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v sslmode=disable password=%v",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresDB,
		c.PostgresPassword,
	)
	conf, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), conf)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}
