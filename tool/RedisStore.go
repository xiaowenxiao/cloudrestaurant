package tool

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const CAPTCHA = "captcha:"

type RedisStore struct {
}

var RedisDb *redis.Client

// 初始化redis配置
func RedisInit() (err error) {
	config := GetConfig().RedisConfig
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Db,
	})
	_, err = RedisDb.Ping().Result()
	return err
}

//实现设置captcha的方法
func (r RedisStore) Set(id string, value string) error {
	key := CAPTCHA + id
	//time.Minute*2：有效时间2分钟
	err := RedisDb.Set(key, value, time.Minute*5).Err()

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
