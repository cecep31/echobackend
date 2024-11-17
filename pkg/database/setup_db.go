package database

import (
	"echobackend/internal/config"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(conf *config.Config) *gorm.DB {
	dsn := conf.GetDSN()

	var config gorm.Config
	if os.Getenv("ENABLE_GORM_LOGGER") != "" {
		config = gorm.Config{}
	} else {
		config = gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn, PreferSimpleProtocol: true,
	}), &config)

	if err != nil {
		log.Fatal(err)
		panic(err.Error()) // Note: This line is redundant after log.Fatal
	}

	return db
}