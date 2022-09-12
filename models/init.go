package models

import (
	"os"
	"treehole_next/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
	},
}

type DBTypeEnum uint

const (
	DBTypeMysql DBTypeEnum = iota
	DBTypeSqlite
)

var DBType DBTypeEnum

func mysqlDB() (*gorm.DB, error) {
	DBType = DBTypeMysql
	return gorm.Open(mysql.Open(config.Config.DbUrl), gormConfig)
}

func sqliteDB() (*gorm.DB, error) {
	DBType = DBTypeSqlite
	err := os.MkdirAll("data", 0750)
	if err != nil {
		panic(err)
	}
	DBType = DBTypeSqlite
	return gorm.Open(sqlite.Open("data/sqlite.db"), gormConfig)
}

func memoryDB() (*gorm.DB, error) {
	DBType = DBTypeSqlite
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), gormConfig)
}

func InitDB() {
	var err error
	switch config.Config.Mode {
	case "production":
		DB, err = mysqlDB()
	case "test":
		DB, err = memoryDB()
		DB = DB.Debug()
	case "bench":
		DB, err = memoryDB()
	case "dev":
		if config.Config.DbUrl == "" {
			DB, err = sqliteDB()
		} else {
			DB, err = mysqlDB()
		}
		DB = DB.Debug()
	default: // sqlite as default
		panic("unknown mode")
	}
	if err != nil {
		panic(err)
	}
	// models must be registered here to migrate into the database
	err = DB.AutoMigrate(
		&Division{},
		&Tag{},
		&Hole{},
		&AnonynameMapping{},
		&Floor{},
		&FloorHistory{},
		&FloorLike{},
		&Report{},
		&UserFavorites{},
	)
	if err != nil {
		panic(err)
	}
}
