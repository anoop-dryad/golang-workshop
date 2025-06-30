package config

import (
	"log"
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
	Env      string
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

	return &cfg, nil
}

func LoadConfig() *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("env", "dev")

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

	return v
}
