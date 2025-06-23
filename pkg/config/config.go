package config

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type TokenConfig struct {
	AccessTokenSecretKey   []byte        `mapstructure:"access_token_secret_key" b64:"true"`
	RefreshTokenPrivateKey []byte        `mapstructure:"refresh_token_private_key" b64:"true"`
	RefreshTokenPublicKey  []byte        `mapstructure:"refresh_token_public_key" b64:"true"`
	AccessTokenDuration    time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration   time.Duration `mapstructure:"refresh_token_duration"`
}

type SecurityConfig struct {
	Token TokenConfig `mapstructure:"token"`
}

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Server   ServerConfig   `mapstructure:"server"`
	Security SecurityConfig `mapstructure:"security"`
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "servling")
	v.SetDefault("database.password", "servling")
	v.SetDefault("database.database", "servling")
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("security.token.access_token_secret_key", base64.StdEncoding.EncodeToString(paseto.NewV4SymmetricKey().ExportBytes()))
	refreshTokenPrivateKey := paseto.NewV4AsymmetricSecretKey()
	v.SetDefault("security.token.refresh_token_private_key", base64.StdEncoding.EncodeToString(refreshTokenPrivateKey.ExportBytes()))
	v.SetDefault("security.token.refresh_token_public_key", base64.StdEncoding.EncodeToString(refreshTokenPrivateKey.Public().ExportBytes()))
	v.SetDefault("security.token.access_token_duration", "15m")
	v.SetDefault("security.token.refresh_token_duration", "24h")
}
func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Config file
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")

	// Environment variables
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read or create config
	if err := v.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err = v.SafeWriteConfigAs("config.toml")
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg, b64decoder()); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func b64decoder() viper.DecoderConfigOption {
	return func(config *mapstructure.DecoderConfig) {
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(config.DecodeHook, decodeB64Hook)
	}
}

func decodeB64Hook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if to.Kind() != reflect.Slice || to.Elem().Kind() != reflect.Uint8 {
		return data, nil
	}
	if from.Kind() == reflect.String {
		return base64.StdEncoding.DecodeString(data.(string))
	}
	if from.Kind() == reflect.Slice && from.Elem().Kind() == reflect.String {
		return base64.StdEncoding.DecodeString(data.([]string)[0])
	}

	return data, nil
}
