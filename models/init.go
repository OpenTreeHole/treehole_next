package models

import (
	"gorm.io/plugin/dbresolver"
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

// Read/Write Splitting
func mysqlDB() (*gorm.DB, error) {
	DBType = DBTypeMysql
	db, err := gorm.Open(mysql.Open(config.Config.DbURL), gormConfig)
	if err != nil {
		return nil, err
	}
	if len(config.Config.MysqlReplicaURLs) == 0 {
		return db, nil
	}
	var replicas []gorm.Dialector
	for _, url := range config.Config.MysqlReplicaURLs {
		replicas = append(replicas, mysql.Open(url))
	}
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(config.Config.DbURL)},
		Replicas: replicas,
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		return nil, err
	}
	return db, nil
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
		if config.Config.DbURL == "" {
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
