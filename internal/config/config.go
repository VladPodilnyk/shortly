package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type StorageConfig struct {
	MongoDbUri        string `mapstructure:"uri"`
	MongoDbName       string `mapstructure:"name"`
	MongoDbCollection string `mapstructure:"collection"`
}

type AppConfig struct {
	Server           ServerConfig  `mapstructure:"server"`
	Storage          StorageConfig `mapstructure:"storage"`
	Environment      string        `mapstructure:"env"`
	Prefix           string        `mapstructure:"prefix"`
	AliasMaxSize     int           `mapstructure:"alias_max_size"`
	RequestPerMinute int           `mapstructure:"requests_per_minute"`
}

type intBuffer []int
type strBuffer []string

func (buf *intBuffer) getFromEnv(variable string) error {
	value := os.Getenv(variable)
	res, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	*buf = append(*buf, res)
	return nil
}

func (buf *strBuffer) getFromEnv(variable string) error {
	value := os.Getenv(variable)
	if len(value) == 0 {
		return fmt.Errorf("couldn't read env variable %s", variable)
	}
	*buf = append(*buf, value)
	return nil
}

func ReadConfig() (AppConfig, error) {
	maybeConfig, err := readConfigFromFile()
	if err == nil {
		return maybeConfig, nil
	}
	fmt.Println("Application config is not provided. Getting config from environment...")
	maybeConfig, err = readConfigFromEnv()
	return maybeConfig, err
}

func readConfigFromFile() (AppConfig, error) {
	configFilepathFlag := flag.String("config", "", "Config file path.")
	flag.Parse()

	if configFilepathFlag == nil || len(*configFilepathFlag) == 0 {
		return AppConfig{}, errors.New("config filepath is not provided")
	}

	filepath := *configFilepathFlag
	fmt.Println("Using the following service config:", filepath)

	configReader := viper.New()
	var config AppConfig

	configReader.SetConfigName(filepath)
	// support only JSON for now
	configReader.SetConfigType("json")
	configReader.AddConfigPath(".")

	err := configReader.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	err = configReader.Unmarshal(&config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, nil
}

// A super lame way to read the service config from environment
func readConfigFromEnv() (AppConfig, error) {
	strValues := []string{"MONGO_DB_URI", "MONGO_DB_NAME", "MONGO_DB_COLLECTION", "PREFIX", "ENVIRONMENT"}
	intValues := []string{"PORT", "ALIAS_MAX_SIZE", "REQUEST_PER_MINUTE"}
	strRes := strBuffer(make([]string, 0, 5))
	intRes := intBuffer(make([]int, 0, 3))

	for _, value := range strValues {
		err := strRes.getFromEnv(value)
		if err != nil {
			return AppConfig{}, err
		}
	}

	for _, value := range intValues {
		err := intRes.getFromEnv(value)
		if err != nil {
			return AppConfig{}, err
		}
	}

	config := AppConfig{
		Storage: StorageConfig{
			MongoDbUri:        strRes[0],
			MongoDbName:       strRes[1],
			MongoDbCollection: strRes[2],
		},
		Prefix:           strRes[3],
		Environment:      strRes[4],
		Server:           ServerConfig{Port: intRes[0]},
		AliasMaxSize:     intRes[1],
		RequestPerMinute: intRes[2],
	}
	return config, nil
}
