package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"sync/atomic"
)

var Config struct {
	Mode          string `env:"MODE" envDefault:"dev"`
	TZ            string `env:"TZ" envDefault:"Asia/Shanghai"`
	Size          int    `env:"SIZE" envDefault:"30"`
	MaxSize       int    `env:"MAX_SIZE" envDefault:"50"`
	TagSize       int    `env:"TAG_SIZE" envDefault:"5"`
	HoleFloorSize int    `env:"HOLE_FLOOR_SIZE" envDefault:"10"`
	Debug         bool   `env:"DEBUG" envDefault:"false"`
	// example: user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=true&loc=Asia%2fShanghai
	// set time_zone in url, otherwise UTC
	// for more detail, see https://github.com/go-sql-driver/mysql#dsn-data-source-name
	DbURL string `env:"DB_URL"`
	// example: MYSQL_REPLICA_URL="db1_dsn,db2_dsn", use ',' as separator
	// should also set time_zone in url
	MysqlReplicaURLs []string `env:"MYSQL_REPLICA_URL"`
	RedisURL         string   `env:"REDIS_URL"` // redis:6379
	NotificationUrl  string   `env:"NOTIFICATION_URL"`
	AuthUrl          string   `env:"AUTH_URL"`
	ElasticsearchUrl string   `env:"ELASTICSEARCH_URL"`
	OpenSearch       bool     `env:"OPEN_SEARCH" envDefault:"true"`
}

var DynamicConfig struct {
	OpenSearch atomic.Bool
}

func init() { // load config from environment variables
	if err := env.Parse(&Config); err != nil {
		panic(err)
	}
	fmt.Println(Config)
	DynamicConfig.OpenSearch.Store(Config.OpenSearch)
}
