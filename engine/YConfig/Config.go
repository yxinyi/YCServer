package YConfig

import "github.com/spf13/viper"

const DEFAULT_PATH = "./config"


func Load(name_ string, core_ interface{}) error {
	_viper := viper.New()
	_viper.SetConfigName(name_)  // name of config file (without extension)
	_viper.SetConfigType("json") // REQUIRED if the config file does not have the extension in the name
	_viper.AddConfigPath(DEFAULT_PATH)  // path to look for the config file in
	_viper.ReadInConfig()
	return _viper.Unmarshal(core_)
}
