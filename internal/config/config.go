package config

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix      = "DONG"
	configFileName = "config.yaml"
)

type DBConfig struct {
	URI  string `mapstructure:"uri"`
	Name string `mapstructure:"name"`
}

type BotConfig struct {
	Token      string `mapstructure:"token"`
	Endpoint   string `mapstructure:"endpoint"`
	Verbose    bool   `mapstructure:"verbose"`
	ListenPort string `mapstructure:"listen_port,omitempty"`
}

type DongConfig struct {
	BotConfig `mapstructure:"bot"`
	DebugMode bool     `mapstructure:"debug,omitempty"`
	DB        DBConfig `mapstructure:"db"`
}

var GlobalConfig *DongConfig

func LoadConfig(cmd *cobra.Command) error {
	GlobalConfig = &DongConfig{}

	// default configs
	viper.SetDefault("bot.listen_port", "80")
	viper.SetDefault("bot.verbose", "false")
	viper.SetDefault("debug", "false")
	viper.SetDefault("db.name", "dong")

	// read config from the config file
	if _, err := os.Stat(configFileName); err == nil {
		viper.SetConfigFile(configFileName)
		logrus.Infof("using config file: %s", viper.ConfigFileUsed())
		err = errors.Wrap(viper.ReadInConfig(), "failed to read config from the config file")
		if err != nil {
			return err
		}
	}

	// read config from the envs
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// bind read configs to the global config
	err := errors.Wrap(
		viper.BindEnv("bot.token", "BOT.TOKEN"), "failed to bind BOT.TOKEN env")
	if err != nil {
		return err
	}
	err = errors.Wrap(
		viper.BindEnv("bot.endpoint", "BOT.ENDPOINT"), "failed to bind BOT_ENDPOINT env")
	if err != nil {
		return err
	}
	err = errors.Wrap(
		viper.BindEnv("bot.verbose", "BOT.VERBOSE"), "failed to bind LISTEN_PORT env")
	if err != nil {
		return err
	}
	err = errors.Wrap(
		viper.BindEnv("db.uri", "DB.URI"), "failed to bind DB.URI env")
	if err != nil {
		return err
	}
	err = errors.Wrap(
		viper.BindEnv("db.name", "DB.NAME"), "failed to bind DB.NAME env")
	if err != nil {
		return err
	}
	err = errors.Wrap(
		viper.BindEnv("bot.listen_port", "BOT.LISTEN_PORT"), "failed to bind BOT.LISTEN_PORT env")
	if err != nil {
		return err
	}
	err = errors.Wrap(
		viper.BindEnv("debug", "DEBUG"), "failed to bind DEBUG_MODE env")
	if err != nil {
		return err
	}

	// bind flags read from the running argument
	err = errors.Wrap(viper.BindPFlags(cmd.Flags()), "failed to bind flags")
	if err != nil {
		return err
	}

	err = errors.Wrap(viper.Unmarshal(&GlobalConfig), "failed to unmarshal the config")
	if err != nil {
		return err
	}
	return nil
}
