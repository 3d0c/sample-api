package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/3d0c/sample-api/pkg/helpers"
)

var db *gorm.DB

func ConnectDatabase() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		helpers.Getenv("DBHOST", "127.0.0.1"),
		helpers.Getenv("DBPORT", "5432"),
		helpers.Getenv("DBUSER", "postgres"),
		helpers.Getenv("DBNAME", "sampleapi"),
	)

	conn, err := gorm.Open("postgres", dsn)
	if err != nil {
		return err
	}

	conn.LogMode(true)

	if err = conn.AutoMigrate(&User{}, &Flight{}).Error; err != nil {
		return err
	}

	db = conn

	return nil
}
