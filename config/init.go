package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/kobeHub/Pegasus-engine/registry"
)

func Init(env string) {
	log.SetFormatter(&log.TextFormatter{})
	// redirect to stdout
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	var config_file string
	if env == "local" {
		config_file = ".env-local"
	} else if env == "test" {
		config_file = ".env-test"
	} else if env == "prod" {
		config_file = ".env-prod"
	} else {
		panic("Error running environment spefic!")
	}
	log.WithFields(log.Fields{
		"file": config_file,
	}).Debug("Initialize config with file")

	viper.SetConfigFile(config_file)
	viper.SetConfigType("dotenv")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.WithFields(log.Fields{
				"file": config_file,
			}).Fatal("Config file not found")
		} else {
			log.Fatalf("Other fatal error: %s", err)
		}
	}

	// init alicloud sdk client
	registry.InitClient()
}
