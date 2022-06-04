package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"treehole_next/config"
)

var DB *gorm.DB

func InitDB() {
	var err error
	if config.Config.Debug {
		DB, err = gorm.Open(sqlite.Open("db/sqlite.db"), &gorm.Config{})
	} else {
		DB, err = gorm.Open(mysql.Open(config.Config.DbUrl), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}
}
