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
	MaxSize int    `default:"10"`
	Debug   bool   `default:"false"`
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
