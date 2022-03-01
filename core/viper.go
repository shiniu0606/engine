package core

import (
	"flag"
	"os"

	"github.com/shiniu0606/engine/core/base"
	"github.com/shiniu0606/engine/core/utils"
	"github.com/spf13/viper"
)

func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" {
			if configEnv := os.Getenv(utils.ConfigEnv); configEnv == "" {
				config = utils.ConfigFile
				base.LogInfo("Use defulat config file,config file path %v\n", utils.ConfigFile)
			} else {
				config = configEnv
				base.LogInfo("Use CONFIG env,config file path %v\n", config)
			}
		} else {
			base.LogInfo("Use command -c param value,config file path %v\n", config)
		}
	} else {
		config = path[0]
		base.LogInfo("Use func Viper() param value,config file path %v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		base.LogFatal("Fatal error config file: %s \n", err)
	}
	return v
}
