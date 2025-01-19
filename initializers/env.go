package initializers

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	// database
	DBHost string `mapstructure:"POSTGRES_HOST"`
	DBUser string `mapstructure:"POSTGRES_USER"`
	DBPass string `mapstructure:"POSTGRES_PASSWORD"`
	DBName string `mapstructure:"POSTGRES_DB"`
	DBPort string `mapstructure:"POSTGRES_PORT"`

	// app
	ServerPort string `mapstructure:"APP_PORT"`

	// JWT Access Token
	AccessTokenPrivate   string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublic    string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenExpiresIn time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`

	// JWT Refresh Token
	RefreshTokenPrivate   string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublic    string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`

	// Redis
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPass string `mapstructure:"REDIS_PASSWORD"`

	// MinIO
	MinioEndpoint      string `mapstructure:"MINIO_ENDPOINT"`
	MinioRootUser      string `mapstructure:"MINIO_ROOT_USER"`
	MinioRootPassword  string `mapstructure:"MINIO_ROOT_PASSWORD"`
	MinioAccessKey     string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey     string `mapstructure:"MINIO_SECRET_KEY"`
	AwsAccessKeyId     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AwsRegion          string `mapstructure:"AWS_REGION"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
