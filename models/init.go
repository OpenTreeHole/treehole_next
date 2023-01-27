package models

import (
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
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
	Logger: logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	),
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

	// set source databases
	source := mysql.Open(config.Config.DbURL)
	db, err := gorm.Open(source, gormConfig)
	if err != nil {
		return nil, err
	}

	// set replica databases
	var replicas []gorm.Dialector
	for _, url := range config.Config.MysqlReplicaURLs {
		replicas = append(replicas, mysql.Open(url))
	}
	err = db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{source},
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

func init() {
	var err error
	switch config.Config.Mode {
	case "production":
		DB, err = mysqlDB()
	case "test":
		fallthrough
	case "bench":
		DB, err = memoryDB()
	case "dev":
		if config.Config.DbURL == "" {
			DB, err = sqliteDB()
		} else {
			DB, err = mysqlDB()
		}
	default:
		panic("unknown mode")
	}
	if err != nil {
		panic(err)
	}

	switch config.Config.Mode {
	case "test":
		fallthrough
	case "dev":
		DB = DB.Debug()
	}

	err = DB.SetupJoinTable(&User{}, "UserFavoriteHoles", &UserFavorite{})
	if err != nil {
		panic(err)
	}

	err = DB.SetupJoinTable(&User{}, "UserLikedFloors", &FloorLike{})
	if err != nil {
		panic(err)
	}

	err = DB.SetupJoinTable(&Hole{}, "Mapping", &AnonynameMapping{})
	if err != nil {
		panic(err)
	}

	// models must be registered here to migrate into the database
	err = DB.AutoMigrate(
		&Division{},
		&Tag{},
		&User{},
		&Floor{},
		&Hole{},
		&Report{},
		&Punishment{},
	)
	if err != nil {
		panic(err)
	}
}
