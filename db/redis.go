package db

import (
	"jumpserver/config"
	"jumpserver/log"
	"time"
	"github.com/garyburd/redigo/redis"
)

var RedisCliPool *redis.Pool

// 初始化redis数据库
func InitRedisPool(c *config.RedisConfig) {
	RedisCliPool = &redis.Pool{
		MaxIdle:     c.MaxIdle,
		MaxActive:   c.MaxActive,
		IdleTimeout: (60 * time.Second),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", c.Host)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH"); err != nil {
				if errclose := c.Close();errclose != nil{
					log.Error("Close redis error by",errclose)
				}
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func CloseRedis(){
	if err := RedisCliPool.Close();err != nil{
		log.Error("Close redis error by",err)
	}
}
