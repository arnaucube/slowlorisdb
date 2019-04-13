package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type Config struct {
	Port string
	Dest string
}

var C Config

func MustRead(c *cli.Context) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")          // adding home directory as first search path
	viper.SetEnvPrefix("slowlorisdb") // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()              // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if c.GlobalString("config") != "" {
		viper.SetConfigFile(c.GlobalString("config"))
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	if err := viper.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
