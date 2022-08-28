package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type MyConfig struct {
	Mode    string `default:"dev"`
	Size    int    `default:"30"`
	MaxSize int    `default:"50"`
	TagSize int    `default:"5"`
	Debug   bool   `default:"false"`
	// example: user:pass@tcp(127.0.0.1:3306)/dbname
	// for more detail, see https://github.com/go-sql-driver/mysql#dsn-data-source-name
	DbUrl    string `default:""`
	MicroUrl string `default:"http://127.0.0.1:8080/api/messages"`
}

var Config MyConfig

func initConfig() { // load config from environment variables
	configType := reflect.TypeOf(Config)
	elem := reflect.ValueOf(&Config).Elem()
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		defaultValue, defaultValueExists := field.Tag.Lookup("default")
		env := os.Getenv(strings.ToUpper(field.Name))
		envExists := env != ""
		if !envExists {
			if !defaultValueExists {
				panic(fmt.Sprintf("Environment variable %s must be set!", field.Name))
			}
			env = defaultValue
		}
		var value interface{}
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
}
