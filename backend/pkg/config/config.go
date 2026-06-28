package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App           AppConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Elasticsearch ElasticsearchConfig
	JWT           JWTConfig
	SMTP          SMTPConfig
	Frontend      FrontendConfig
	RateLimit     RateLimitConfig
	Scheduler     SchedulerConfig
}

type AppConfig struct {
	Env  string
	Port string
	Name string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type ElasticsearchConfig struct {
	Addresses []string
	Username  string
	Password  string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

type FrontendConfig struct {
	URL string
}

type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

type SchedulerConfig struct {
	JobIntervalHours int
	DigestCron       string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	}

	esAddresses := viper.GetString("ES_ADDRESSES")
	var esAddressList []string
	if esAddresses != "" {
		esAddressList = strings.Split(esAddresses, ",")
	} else {
		esAddressList = []string{"http://localhost:9200"}
	}

	return &Config{
		App: AppConfig{
			Env:  getEnvOrDefault("APP_ENV", "development"),
			Port: getEnvOrDefault("APP_PORT", "8080"),
			Name: getEnvOrDefault("APP_NAME", "CareerCopilot"),
		},
		Database: DatabaseConfig{
			Host:            getEnvOrDefault("DB_HOST", "localhost"),
			Port:            getEnvOrDefault("DB_PORT", "5432"),
			User:            getEnvOrDefault("DB_USER", "careercopilot"),
			Password:        getEnvOrDefault("DB_PASSWORD", "careercopilot_secret"),
			Name:            getEnvOrDefault("DB_NAME", "careercopilot"),
			SSLMode:         getEnvOrDefault("DB_SSLMODE", "disable"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: time.Duration(viper.GetInt("DB_CONN_MAX_LIFETIME")) * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnvOrDefault("REDIS_HOST", "localhost"),
			Port:     getEnvOrDefault("REDIS_PORT", "6379"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		Elasticsearch: ElasticsearchConfig{
			Addresses: esAddressList,
			Username:  viper.GetString("ES_USERNAME"),
			Password:  viper.GetString("ES_PASSWORD"),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnvOrDefault("JWT_ACCESS_SECRET", "default_access_secret"),
			RefreshSecret: getEnvOrDefault("JWT_REFRESH_SECRET", "default_refresh_secret"),
			AccessExpiry:  time.Duration(viper.GetInt("JWT_ACCESS_EXPIRY")) * time.Minute,
			RefreshExpiry: time.Duration(viper.GetInt("JWT_REFRESH_EXPIRY")) * time.Minute,
		},
		SMTP: SMTPConfig{
			Host:      getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			Port:      viper.GetInt("SMTP_PORT"),
			Username:  viper.GetString("SMTP_USERNAME"),
			Password:  viper.GetString("SMTP_PASSWORD"),
			FromName:  getEnvOrDefault("SMTP_FROM_NAME", "CareerCopilot"),
			FromEmail: getEnvOrDefault("SMTP_FROM_EMAIL", "noreply@careercopilot.io"),
		},
		Frontend: FrontendConfig{
			URL: getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
		},
		RateLimit: RateLimitConfig{
			Requests: viper.GetInt("RATE_LIMIT_REQUESTS"),
			Duration: time.Duration(viper.GetInt("RATE_LIMIT_DURATION")) * time.Second,
		},
		Scheduler: SchedulerConfig{
			JobIntervalHours: viper.GetInt("SCHEDULER_JOB_INTERVAL_HOURS"),
			DigestCron:       getEnvOrDefault("SCHEDULER_DIGEST_CRON", "0 7 * * *"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	val := viper.GetString(key)
	if val == "" {
		return defaultValue
	}
	return val
}
