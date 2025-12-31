package database

import (
	"errors"
	"math/rand"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/ryansiau/KeepUpdated/go/model"
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

type ConnectionTest struct {
	ID int `gorm:"primaryKey"`
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

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Content{},
		&ConnectionTest{},
	)
	if err != nil {
		panic(err)
	}

	return nil
}

func CheckConnectionCapability(db *gorm.DB) error {
	var temp ConnectionTest

	var rollRandomID = true
	var randomId int

	for rollRandomID {
		randomId = int(rand.Int63n(1<<31 - 1))

		res := db.Find(&temp, randomId)
		if res.Error != nil {
			return res.Error
		}
		if temp.ID != randomId {
			rollRandomID = false
		}
	}

	res := db.Create(&ConnectionTest{ID: randomId})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 1 {
		return errors.New("failed to create record")
	}

	res = db.Delete(&ConnectionTest{}, randomId)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected != 1 {
		return errors.New("failed to delete record")
	}

	return nil
}
