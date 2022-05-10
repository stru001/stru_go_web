package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"stru_web/settings"
)

var rdb *redis.Client
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d",cfg.Host,cfg.Port), // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		zap.L().Error("connect redis fail.",zap.Error(err))
		return err
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}