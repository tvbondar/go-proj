// Загрузка настроек (из env или YAML)
// Этот файл вынесет конфигурации (DSN для PostgreSQL, Kafka адреса, HTTP порт)
package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBDSN     string `mapstructure:"db_dsn"`
	KafkaAddr string `mapstructure:"kafka_addr"`
	HTTPPort  string `mapstructure:"http_port"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // Имя файла конфига (config.yaml)
	viper.SetConfigType("yaml")   // Тип файла
	viper.AddConfigPath(".")      // Путь к файлу

	// Поддержка переменных окружения
	viper.AutomaticEnv()

	// Установка значений по умолчанию
	if err := viper.BindEnv("kafka_addr", "KAFKA_BROKERS"); err != nil {
		return nil, fmt.Errorf("failed to bind KAFKA_BROKERS env: %w", err)
	}
	viper.SetDefault("http_port", ":8081")

	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if cfg.DBDSN == "" {
		if viper.GetString("DB_DSN") != "" {
			cfg.DBDSN = viper.GetString("DB_DSN")
		} else {
			return nil, errors.New("db_dsn is not set; set DB_DSN env or config.yaml")
		}
	}
	return &cfg, nil
}
