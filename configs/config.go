package configs

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string
	HTTPPort string
	KeyPath
	Database
	GCP
}

type KeyPath struct {
	PublicKey  string
	PrivateKey string
}

type Database struct {
	PostgresDB struct {
		Host     string
		Port     string
		Username string
		Password string
		DbName   string
		SSLMode  bool
	}
}

type GCP struct {
	ProjectID      string
	ServiceAccount string
}

var config Config

func InitViper(path string, env string) {
	switch env {
	case "local":
		viper.SetConfigName("local-config")
	case "develop":
		viper.SetConfigName("develop-config")
	case "staging":
		viper.SetConfigName("staging-config")
	default:
		viper.SetConfigName("config")
	}
	log.Println("Running on environment:", env)
	viper.AddConfigPath(path)
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}

func GetViper() *Config {
	return &config
}
