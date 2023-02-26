package tool

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

var (
	RedisDb *redis.Client
)

// 创建redis连接
func InitRedis() {
	config := GetConfig().RedisConfig
	RedisDb := redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Db,
	})
	_, err := RedisDb.Ping().Result()
	if err != nil {
		//连接失败
		println(err)
	}
	log.Println("Redis连接成功")

}

const CAPTCHA = "captcha:"

type RedisStore struct {
}

//实现设置captcha的方法
func (r RedisStore) Set(id string, value string) error {
	config := GetConfig().RedisConfig
	RedisDb := redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Db,
	})
	_, erra := RedisDb.Ping().Result()
	if erra != nil {
		//连接失败
		println(erra)
	}
	log.Println("Redis连接成功")
	key := CAPTCHA + id
	//time.Minute*2：有效时间2分钟
	err := RedisDb.Set(key, value, time.Minute*2).Err()

	return err
}

//实现获取captcha的方法
func (r RedisStore) Get(id string, clear bool) string {

	key := CAPTCHA + id
	val, err := RedisDb.Get(key).Result()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if clear {
		//clear为true，验证通过，删除这个验证码
		err := RedisDb.Del(key).Err()
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}
	return val
}

//实现验证captcha的方法
func (r RedisStore) Verify(id, answer string, clear bool) bool {

	v := RedisStore{}.Get(id, clear)
	//fmt.Println("key:"+id+";value:"+v+";answer:"+answer)
	return v == answer
}
