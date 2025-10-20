// Загрузка настроек (из env или YAML)
// Этот файл вынесет конфигурации (DSN для PostgreSQL, Kafka адреса, HTTP порт)
// из main.go в отдельный модуль с использованием библиотеки viper
package config

import (
	"errors"

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

	// Поддержка переменных окружения
	viper.AutomaticEnv()

	// Установка значений по умолчанию (без секретов!)
	viper.SetDefault("kafka_addr", "localhost:9092")
	viper.SetDefault("http_port", ":8081")
	// По безопасности: не ставим пароль/DSN по-умолчанию. Требуем явно задать DB_DSN или config.yaml.
	// viper.SetDefault("db_dsn", "")

	_ = viper.ReadInConfig() // если файла нет — используем env/дефолты

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Проверка обязательных значений
	if cfg.DBDSN == "" {
		// Попробуем получить из окружения DB_DSN
		if viper.GetString("DB_DSN") != "" {
			cfg.DBDSN = viper.GetString("DB_DSN")
		} else {
			return nil, errors.New("db_dsn is not set; set DB_DSN env or config.yaml")
		}
	}
	return &cfg, nil
}
