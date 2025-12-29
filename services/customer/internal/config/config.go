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
	Auth     AuthConfig
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type GRPCConfig struct {
	Port string `mapstructure:"port"`
}

type AuthConfig struct {
	// jwt settings
	JWTSecret      string `mapstructure:"jwt_secret"`
	JWTExpiryHours int    `mapstructure:"jwt_expiry_hours"`

	// password settings
	BcryptCost        int `mapstructure:"bcrypt_cost"`
	MinPasswordLength int `mapstructure:"min_password_length"`

	// security settings
	LockoutThreshold int `mapstructure:"lockout_treshold"`
	LockoutDuration  int `mapstructure:"lockout_duration"`

	// email verification settings
	VerificationTokenExpiry int `mapstructure:"verification_token_expiry"`

	// password reset
	ResetTokenExpiry int `mapstructure:"reset_token_expiry"`
}

func LoadConfig() (*Config, error) {
	// load the env file
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	// print the env file here
	envFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(envFile); err != nil {
		godotenv.Load()
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")

	// set deafult
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("grpc.port", 50051)
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.ssl_mode", "disable")

	// auth defaults
	viper.SetDefault("auth.jwt_secret", "value-to-be-changed-in-prod")
	viper.SetDefault("auth.jwt_expiry_hours", 24)
	viper.SetDefault("auth.bcrypt_cost", 12)
	viper.SetDefault("auth.min_password_length", 8)
	viper.SetDefault("auth.lockout_threshold", 5)
	viper.SetDefault("auth.lockout_duration_minutes", 15)
	viper.SetDefault("auth.reset_token_expiry_minutes", 60)

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
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	return &config, nil
}

func bindEnvVars() {
	// server variables
	viper.BindEnv("server.host", "CUSTOMER_SERVER_HOST")
	viper.BindEnv("server.port", "CUSTOMER_SERVER_PORT")

	// grpc variables
	viper.BindEnv("grpc.port", "CUSTOMER_GRPC_PORT")

	// database variables
	viper.BindEnv("database.host", "CUSTOMER_DB_HOST")
	viper.BindEnv("database.port", "CUSTOMER_DB_PORT")
	viper.BindEnv("database.user", "CUSTOMER_DB_USER")
	viper.BindEnv("database.password", "CUSTOMER_DB_PASSWORD")
	viper.BindEnv("database.name", "CUSTOMER_DB_NAME")

	// auth variables
	viper.BindEnv("auth.jwt_secret", "CUSTOMER_AUTH_JWT_SECRET")
	viper.BindEnv("auth.bcrypt_cost", "CUSTOMER_AUTH_BCRYPT_COST")
	viper.BindEnv("auth.jwt_expiry_hours", "CUSTOMER_AUTH_JWT_EXPIRY_HOURS")
	viper.BindEnv("auth.min_password_length", "CUSTOMER_AUTH_MIN_PASSWORD_LENGTH")
	viper.BindEnv("auth.lockout_threshold", "CUSTOMER_AUTH_LOCKOUT_THRESHOLD")
	viper.BindEnv("auth.lockout_duration_minutes", "CUSTOMER_AUTH_LOCKOUT_DURATION_MINUTES")
	viper.BindEnv("auth.reset_token_expiry_minutes", "CUSTOMER_AUTH_RESET_TOKEN_EXPIRY_MINUTES")
	viper.BindEnv("auth.verification_token_expiry_hours", "CUSTOMER_AUTH_VERIFICATION_EXPIRY_HOURS")

}

func (c *Config) Validate() error {
	// check jwt secret[should be changed in production]
	if c.Auth.JWTSecret == "value-to-be-changed-in-prod" {
		log.Println("WARNING: using default jwt secret. CHANGE THIS IN PRODUCTION!")
	}

	if len(c.Auth.JWTSecret) < 32 {
		log.Println("WARNING: consider using a longer secret")
	}

	// check bcrypt cost
	if c.Auth.BcryptCost < 4 || c.Auth.BcryptCost > 31 {
		return fmt.Errorf("bcrypt cost must be between 4 and 31, got %d:", c.Auth.BcryptCost)
	}

	// check jwt expiry
	if c.Auth.JWTExpiryHours < 1 {
		return fmt.Errorf("JWT expiry hours must be at least 1, got %d", c.Auth.JWTExpiryHours)
	}

	// check password length
	if c.Auth.MinPasswordLength < 8 {
		return fmt.Errorf("minimum password length must be at least 8, got %d", c.Auth.MinPasswordLength)
	}

	// check lockout settings
	if c.Auth.LockoutThreshold < 1 {
		return fmt.Errorf("lockout threshold must be at least 1, got %d", c.Auth.LockoutThreshold)
	}

	if c.Auth.LockoutDuration < 1 {
		return fmt.Errorf("lockout duration must be at least 1 minute, got %d", c.Auth.LockoutDuration)
	}

	return nil
}
