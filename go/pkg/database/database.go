package database

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ryansiau/utilities/go/model"
)

type Config struct {
	DatabaseType string `yaml:"type"`

	// SQLite
	Filepath string `yaml:"filepath"`

	// PostgreSQL & MySQL
	DatabaseName string `yaml:"database_name"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
}

func NewDB(conf *Config) (*gorm.DB, error) {
	switch conf.DatabaseType {
	case "sqlite":
		gormDB, err := gorm.Open(sqlite.Open(conf.Filepath), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		return gormDB, nil
	default:
		return nil, errors.New("Unknown database type")
	}
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Content{},
	)
	if err != nil {
		panic(err)
	}

	return nil
}
