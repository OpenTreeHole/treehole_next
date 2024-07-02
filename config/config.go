package config

import (
	"sync/atomic"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v9"
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
	MysqlReplicaURLs   []string `env:"MYSQL_REPLICA_URL"`
	RedisURL           string   `env:"REDIS_URL"` // redis:6379
	NotificationUrl    string   `env:"NOTIFICATION_URL"`
	MessagePurgeDays   int      `envDefault:"7" env:"MESSAGE_PURGE_DAYS"`
	AuthUrl            string   `env:"AUTH_URL"`
	ElasticsearchUrl   string   `env:"ELASTICSEARCH_URL"`
	OpenSearch         bool     `env:"OPEN_SEARCH" envDefault:"true"`
	OpenFuzzName       bool     `env:"OPEN_FUZZ_NAME" envDefault:"false"`
	UserAllShowHidden  bool     `env:"USER_ALL_HIDDEN" envDefault:"false"`
	AdminOnly          bool     `env:"ADMIN_ONLY" envDefault:"false"`
	HolePurgeDivisions []int    `env:"HOLE_PURGE_DIVISIONS" envDefault:"2"`
	HolePurgeDays      int      `env:"HOLE_PURGE_DAYS" envDefault:"30"`
	OpenSensitiveCheck bool     `env:"OPEN_SENSITIVE_CHECK" envDefault:"true"`

	YiDunBusinessIdText  string   `env:"YI_DUN_BUSINESS_ID_TEXT" envDefault:""`
	YiDunBusinessIdImage string   `env:"YI_DUN_BUSINESS_ID_IMAGE" envDefault:""`
	YiDunSecretId        string   `env:"YI_DUN_SECRET_ID" envDefault:""`
	YiDunSecretKey       string   `env:"YI_DUN_SECRET_KEY" envDefault:""`
	ValidImageUrl        []string `env:"VALID_IMAGE_URL"`
	UrlHostnameWhitelist []string `env:"URL_HOSTNAME_WHITELIST"`
	ExternalImageHost    string   `env:"EXTERNAL_IMAGE_HOSTNAME" envDefault:""`
}

var DynamicConfig struct {
	OpenSearch atomic.Bool
}

func InitConfig() { // load config from environment variables
	if err := env.Parse(&Config); err != nil {
		log.Fatal().Err(err).Send()
	}
	log.Info().Any("config", Config).Msg("init config")
	DynamicConfig.OpenSearch.Store(Config.OpenSearch)
}
