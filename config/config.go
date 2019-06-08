package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

type Config struct {
	StoragePath string
	Port        int
	Dest        string
	AuthNodes   []string // PubKs in hex format of the AuthNodes for the blockchain
}

func MustRead(c *cli.Context) (*Config, error) {
	var config Config

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
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
