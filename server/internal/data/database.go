package data

import (
	"fmt"

	"github.com/fressive/pocman/server/internal/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Initialize the database from server config. Ensure that the
// config is loaded before calling this function.
func InitDatabase() error {
	var err error

	dbconf := conf.ServerConfig.Data.Database

	switch dbconf.Driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(dbconf.Source))
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dbconf.Source), &gorm.Config{})
	case "postgres":
		DB, err = gorm.Open(postgres.Open(dbconf.Source), &gorm.Config{})
	default:
		err = fmt.Errorf("unknown database driver %s", conf.ServerConfig.Data.Database.Driver)
	}

	return err
}
