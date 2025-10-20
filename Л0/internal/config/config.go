// Загрузка настроек (из env или YAML)
// Этот файл вынесет конфигурации (DSN для PostgreSQL, Kafka адреса, HTTP порт)
// из main.go в отдельный модуль с использованием библиотеки viper
package config

import (
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
	viper.AddConfigPath(".")      // Путь к файлу (в корне проекта)

	// Установка значений по умолчанию
	viper.SetDefault("db_dsn", "user=postgres password=pass dbname=orders_db host=postgres port=5432 sslmode=disable")
	viper.SetDefault("kafka_addr", "kafka:9092")
	viper.SetDefault("http_port", ":8081")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err // Ошибка парсинга файла
		}
		// Если файла нет, используем значения по умолчанию
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
