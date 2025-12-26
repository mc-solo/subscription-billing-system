package config

import (
	"fmt"
	"os"

	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	GRPC     GRPCConfig
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"host"`
	User     string `mapstructure:"host"`
	Password string `mapstructure:"host"`
	Name     string `mapstructure:"host"`
	SSLMode  string `mapstructure:"host"`
}

type GRPCConfig struct {
	Port string `mapstructure:"host"`
}

func LoadConfig() (*Config, error) {
	// load the env file

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	envFile := fmt.Sprintf(".env.%s", env)

	if err := godotenv.Load(envFile); err != nil {
		godotenv.Load()
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// set deafults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("grpc.port", 50051)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.ssl_mode", "disable")

	// read the config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found, using default and environment variables")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// bind env vars
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CUSTOMER")

	// bind specific env vars
	bindEnvVars()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return &config, nil
}

func bindEnvVars() {
	viper.BindEnv("server.host", "CUSTOMER_SERVER_HOST")
	viper.BindEnv("server.port", "CUSTOMER_SERVER_PORT")
	viper.BindEnv("grpc.port", "CUSTOMER_GRPC_PORT")
	viper.BindEnv("database.host", "CUSTOMER_DB_HOST")
	viper.BindEnv("database.port", "CUSTOMER_DB_PORT")
	viper.BindEnv("database.user", "CUSTOMER_DB_USER")
	viper.BindEnv("database.password", "CUSTOMER_DB_PASSWORD")
	viper.BindEnv("database.name", "CUSTOMER_DB_NAME")
}
