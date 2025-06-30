package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Cors     CorsConfig
	Logger   LoggerConfigs
	Env      EnvConfig
	Mqtt     []MqttConfig
	Kinesis  []KinesisConfig
	Sqs      SqsConfig
	Jwks     JwksCofig
}

type EnvConfig struct {
	Stage   string
	AppName string
}

type KinesisConfig struct {
	Name string
}

type JwksCofig struct {
	AuthEndpoint string
}

type MqttConfig struct {
	Id        string
	BrokerUrl string
	Topic     string
}

type SqsConfig struct {
	Name string
}

type ServerConfig struct {
	Port    string
	RunMode string
	Domain  string
}

type LoggerConfigs struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

type PostgresConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	Db                 string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleCheckFrequency time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
}

type CorsConfig struct {
	AllowOrigins string
}

func GetConfig() *Config {
	cfg, err := ParseConfig(LoadConfig())
	if err != nil {
		log.Fatalf("Error in parse config %v", err)
	}

	return cfg
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var cfg Config
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}

	cfg.Kinesis = loadKinesisStreams()
	mqttClients, err := loadMQTTClients()
	if err != nil {
		log.Fatalf("mqtt client parsing failed : %s", err.Error())
	}
	cfg.Mqtt = mqttClients

	return &cfg, nil
}

func LoadConfig() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Server config
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.runmode", "debug")
	v.SetDefault("server.domain", "localhost")

	// Logger config
	v.SetDefault("logger.filepath", "logs/")
	v.SetDefault("logger.encoding", "json")
	v.SetDefault("logger.level", "debug")
	v.SetDefault("logger.logger", "zap")

	// Postgres config
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", "5432")
	v.SetDefault("postgres.user", "postgres")
	v.SetDefault("postgres.password", "postgres")
	v.SetDefault("postgres.dbname", "postgres")
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.maxidleconns", 15)
	v.SetDefault("postgres.maxopenconns", 100)
	v.SetDefault("postgres.connmaxlifetime", 5)

	// Redis config
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.password", "password")
	v.SetDefault("redis.db", "0")
	v.SetDefault("redis.dialtimeout", 5)
	v.SetDefault("redis.readtimeout", 5)
	v.SetDefault("redis.writetimeout", 5)
	v.SetDefault("redis.idlecheckfrequency", 500)
	v.SetDefault("redis.poolsize", 10)
	v.SetDefault("redis.pooltimeout", 15)

	// Cors config
	v.SetDefault("cors.alloworigins", "*")

	// Env config
	v.SetDefault("env.stage", "dev")
	v.SetDefault("env.appname", "api")

	// Sqs config
	v.SetDefault("sqs.name", "")

	// Jwks config
	v.SetDefault("jwks.authendpoint", "")

	return v
}

func loadKinesisStreams() []KinesisConfig {
	var kinesisConfigs []KinesisConfig

	for i := 0; ; i++ {
		envVar := fmt.Sprintf("KINESIS_STREAM_%d", i)
		envValue := os.Getenv(envVar)
		if envValue == "" {
			break
		}

		kinesisConfig := KinesisConfig{
			Name: envValue,
		}
		kinesisConfigs = append(kinesisConfigs, kinesisConfig)
	}

	return kinesisConfigs
}

func loadMQTTClients() ([]MqttConfig, error) {
	var mqttClients []MqttConfig

	for i := 0; ; i++ {
		envVar := fmt.Sprintf("MQTT_CLIENT_%d", i)
		envVal := os.Getenv(envVar)
		if envVal == "" {
			break
		}

		var cfg MqttConfig
		if err := json.Unmarshal([]byte(envVal), &cfg); err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", envVar, err)
		}

		mqttClients = append(mqttClients, cfg)
	}

	specialClients := []string{"MQTT_CLIENT_GTW"}
	for _, envVar := range specialClients {
		envVal := os.Getenv(envVar)
		if envVal == "" {
			break
		}

		var cfg MqttConfig
		if err := json.Unmarshal([]byte(envVal), &cfg); err != nil {
			return nil, fmt.Errorf("error parsing %s: %w", envVar, err)
		}

		mqttClients = append(mqttClients, cfg)
	}

	return mqttClients, nil
}
