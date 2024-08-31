package config

import (
	"log"
	"reflect"

	"github.com/spf13/viper"
)

// Set default environments
var Env = struct {
	Env    string `mapstructure:"ENV"`
	Port   string `mapstructure:"PORT"`
	Apikey string `mapstructure:"APIKEY"`
	Cors   string `mapstructure:"CORS"`

	DBURI  string `mapstructure:"DB_URI"`
	DBName string `mapstructure:"DB_NAME"`

	ElasticHost  string `mapstructure:"ELASTIC_HOST"`
	ElasticIndex string `mapstructure:"ELASTIC_INDEX"`

	KafkaBrokers       string `mapstructure:"KAFKA_BROKERS"`
	KafkaVersion       string `mapstructure:"KAFKA_VERSION"`
	KafkaConsumerGroup string `mapstructure:"KAFKA_CONSUMER_GROUP"`
}{
	Env:                "production",
	Port:               "3000",
	Cors:               "*",
	KafkaVersion:       "2.1.1",
	KafkaConsumerGroup: "example_group",
}

func NewAppInitEnvironment() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Environment variables not used from .env")
	}

	// try load settings from env vars
	r := reflect.TypeOf(Env)
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i).Tag.Get("mapstructure")
		viper.BindEnv(f)
	}

	if err := viper.Unmarshal(&Env); err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}
}
