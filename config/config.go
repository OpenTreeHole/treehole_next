package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type MyConfig struct {
	Mode    string `default:"dev" env:"MODE"`
	Size    int    `default:"10" env:"SIZE"`
	MaxSize int    `default:"30" env:"MAX_SIZE"`
	TagSize int    `default:"5" env:"TAG_SIZE"`
	Debug   bool   `default:"false" env:"DEBUG"`
	// example: user:pass@tcp(127.0.0.1:3306)/dbname
	// for more detail, see https://github.com/go-sql-driver/mysql#dsn-data-source-name
	DBURL    string `default:"" env:"DB_URL"`
	RedisURL string `default:"redis:6379" env:"REDIS_URL"`
}

var Config MyConfig

func initConfig() { // load config from environment variables
	configType := reflect.TypeOf(Config)
	elem := reflect.ValueOf(&Config).Elem()
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		// get default value
		defaultValue, defaultValueExists := field.Tag.Lookup("default")
		// get env variable name
		envName, ok := field.Tag.Lookup("env")
		if !ok {
			envName = strings.ToUpper(field.Name)
		}
		// get env variable value
		env := os.Getenv(envName)
		envExists := env != ""
		if !envExists {
			if !defaultValueExists {
				panic(fmt.Sprintf("Environment variable %s must be set!", field.Name))
			}
			env = defaultValue
		}
		var value any
		var err error
		switch field.Type.Kind() {
		case reflect.String:
			value = env
		case reflect.Int:
			value, err = strconv.Atoi(env)
			if err != nil {
				panic(fmt.Sprintf("Environment variable %s must be an int!", field.Name))
			}
		case reflect.Bool:
			lower := strings.ToLower(env)
			if lower == "true" {
				value = true
			} else if lower == "false" {
				value = false
			} else {
				panic(fmt.Sprintf("Environment variable %s must be a boolean!", field.Name))
			}
		default:
			panic("Now only supports string, int and bool as config")
		}
		elem.FieldByName(field.Name).Set(reflect.ValueOf(value))
	}
}

func InitConfig() {
	initConfig()
	if Config.Mode != "production" {
		Config.Debug = true
	}
	initCache()
}
