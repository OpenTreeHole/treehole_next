package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"treehole_next/config"
)

var DB *gorm.DB
var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	},
}

func InitDB() {
	var err error
	if config.Config.Debug {
		DB, err = gorm.Open(sqlite.Open("db/sqlite.db"), gormConfig)
	} else {
		DB, err = gorm.Open(mysql.Open(config.Config.DbUrl), gormConfig)
	}
	if err != nil {
		panic(err)
	}
}
