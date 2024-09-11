package global

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GlobalConfig struct {
	Verbose        bool            `mapstructure:"verbose"`
	APIConfig      *APIConfig      `mapstructure:"api"`
	PostgresConfig *PostgresConfig `mapstructure:"pg"`
}

type APIConfig struct {
	Bind string `mapstructure:"bind"`
	Port int    `mapstructure:"port"`
}

type PostgresConfig struct {
	Hostname string  `mapstructure:"hostname"`
	Database string  `mapstructure:"db"`
	SSL      string  `mapstructure:"sslmode"`
	Port     int     `mapstructure:"port"`
	Runtime  *PGAuth `mapstructure:"runtime"`
	Manage   *PGAuth `mapstructure:"manage"`
}

func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s dbname=%s", p.Hostname, p.Runtime.Username, p.Runtime.Password, p.Port, p.SSL, p.Database)
}

type PGAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

//type PostgresConfig struct {
//	Hostname string `mapstructure:"hostname"`
//	Username string `mapstructure:"username"`
//	Password string `mapstructure:"password"`
//	Schema   string `mapstructure:"schema"`
//	SSL      string `mapstructure:"sslmode"`
//	Port     int    `mapstructure:"port"`
//}
//
//func (p *PostgresConfig) DSN(manage bool) string {
//	if manage {
//		return fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s", p.Hostname, p.Username, p.Password, p.Port, p.SSL)
//	}
//	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s", p.Hostname, p.Username, p.Password, p.Schema, p.Port, p.SSL)
//}

func DefaultConfig() *GlobalConfig {
	return &GlobalConfig{}
}

func InitConfig(pathOverride string) error {
	logger := zap.S()
	defer logger.Sync()

	if pathOverride != "" {
		viper.SetConfigFile(pathOverride)
	} else {
		viper.AddConfigPath("./")
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", App))
		viper.SetConfigType("yaml")
		viper.SetConfigName(App)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorw("failed to read initial config", "err", err)
		return errors.New("no config file found")
	}

	cfg := DefaultConfig()
	err := viper.Unmarshal(cfg)
	if err != nil {
		logger.Errorw("failed to parse config file", "path", viper.ConfigFileUsed(), "err", err)
		return errors.New("invalid config file")
	}
	// TODO: reinstate validation once you get there and all lol
	//	if err := cfg.Validate(); err != nil {
	//		return err
	//	}

	_gcMu.Lock()
	_gC = cfg
	_gcMu.Unlock()

	logger.Debugw("config loaded", "path", viper.ConfigFileUsed())
	return nil
}
