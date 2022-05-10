package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AppConfig struct {
	Name 	string `mapstructure:"name"`
	Mode 	string `mapstructure:"mode"`
	Port 	int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
	*LogConfig     `mapstructure:"log"`
	*MySQLConfig   `mapstructure:"mysql"`
	*RedisConfig   `mapstructure:"redis"`
}
type LogConfig struct {
	Level 		string `mapstructure:"level"`
	FileName 	string `mapstructure:"filename"`
	MaxSize		int    `mapstructure:"max_size"`
	MaxAge 		int    `mapstructure:"max_age"`
	MaxBackups 	int    `mapstructure:"max_backups"`
}
type MySQLConfig struct {
	Host 			string `mapstructure:"host"`
	Port 			int    `mapstructure:"port"`
	DbName 			string `mapstructure:"dbname"`
	User 			string `mapstructure:"user"`
	Password 		string `mapstructure:"password"`
	MaxOpenConns 	int    `mapstructure:"max_open_conns"`
	MaxIdleConns 	int    `mapstructure:"max_idle_conns"`
	MaxLiftTime		int    `mapstructure:"max_lift_time"`
}
type RedisConfig struct {
	Host 			string `mapstructure:"host"`
	Port 			int    `mapstructure:"port"`
	Password 		string `mapstructure:"password"`
}

var Conf = new(AppConfig)

func Init(filePath string) (err error) {
	//viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.SetConfigFile(filePath)
	err = viper.ReadInConfig()
	if err != nil {
		zap.L().Error("init viper config fail",zap.Error(err))
		return err
	}
	if err = viper.Unmarshal(Conf); err != nil {
		zap.L().Error("viper.Unmarshal fail",zap.Error(err))
		return err
	}

	// 支持配置文件变化时监听
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.L().Info("config change...")
		if err = viper.Unmarshal(Conf); err != nil {
			zap.L().Error("viper.Unmarshal config change fail",zap.Error(err))
		}
	})
	return nil
}
