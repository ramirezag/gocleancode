package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
	"os"
	"reflect"
	"strconv"
)

type Configuration struct {
	UploadDir string `env:"UPLOAD_DIR"`
	Host      string `env:"HOST"`
	Port      int    `env:"APP_PORT"`
	DbHost    string `env:"DB_HOST"` // Defaults to localhost
	DbPort    int    `env:"DB_PORT"` // Defaults to 3306
	DbUser    string `env:"DB_USER"`
	DbPass    string `env:"DB_PASS"`
	DbName    string `env:"DB_NAME"`
}

func New() Configuration {
	config := Configuration{}
	configPath := getConfigPath()
	err := gonfig.GetConf(configPath, &config)
	if err != nil {
		log.Warn("Failed to fill config from ", configPath)
	}
	fillFromEnvVariables(&config)
	if config.Port == 0 {
		config.Port = 8000
	}
	if config.DbHost == "" {
		config.DbHost = "127.0.0.1"
	}
	if config.DbPort == 0 {
		config.DbPort = 3306
	}
	return config
}

func getConfigPath() string {
	env := os.Getenv("ENV")
	if len(env) == 0 {
		env = "development"
	}
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configPath := workingDir + "/" + "config." + env + ".json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		configPath = workingDir + "/" + "config.default.json"
	}
	return configPath
}

func fillFromEnvVariables(config *Configuration) {
	// Due to gonfig bug, let's fill the values from environment variables manually
	typ := reflect.TypeOf(*config)
	elem := reflect.ValueOf(config).Elem()
	for i := 0; i < typ.NumField(); i++ {
		fT := typ.Field(i)
		f := elem.Field(i)
		if f.IsValid() {
			envVal := os.Getenv(fT.Tag.Get("env"))
			if envVal != "" {
				kind := f.Kind()
				if kind == reflect.Int || kind == reflect.Int64 {
					setStringToInt(f, envVal, 64)
				} else {
					f.SetString(envVal)
				}
			}
		}
	}
}

func setStringToInt(f reflect.Value, value string, bitSize int) {
	convertedValue, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		if !f.OverflowInt(convertedValue) {
			f.SetInt(convertedValue)
		}
	}
}
