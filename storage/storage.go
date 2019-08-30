package storage

import (
	"delay_queue_callback/config"
	"errors"
	"github.com/go-redis/redis"
)

var Db Storage

type Storage struct {
	client *redis.Client
}

var client *redis.Client

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Conf.RedisHost,
		Password: config.Conf.RedisPasswd,
	})

	client = redisClient
}

func Get(key string) (string, error) {
	val, err := client.Get(key).Result()
	if err != nil {
		return "", err
	}

	return val, err

}

func Set(key string, value string) (error) {
	err := client.Set(key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func Del(key string) (error) {
	err := client.Del(key).Err()
	if err != nil {
		return err
	}

	return nil
}

//普通队列推入
func Push(key string, value string) (error) {
	err := client.LPush(key, value).Err()
	if err != nil {
		return err
	}

	return nil
}

//普通队列弹出
func Pop(key string, value string) (string, error) {
	val, err := client.RPop(key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

//有序队列增加
func Zadd(key string, value interface{}, score int64) error {
	res, err := client.ZAdd(key, redis.Z{float64(score), value}).Result()

	if res == 0 {
		return errors.New("加入有序队列失败")
	}

	if err != nil {
		return err
	}

	return nil
}

//从有序队列中删除
func Zrem(key string, member interface{}) error {
	err := client.ZRem(key, member).Err()

	if err != nil {
		return err
	}

	return nil
}

//从有序队列取出排名最靠前的

func GetTop(key string) (string, int64, error) {
	val, err := client.ZRangeWithScores(key, 0, 0).Result();
	if err != nil {
		return "", 0, err
	}

	if len(val) == 0 {
		return "", 0, nil
	}

	score := int64(val[0].Score)
	value := val[0].Member.(string)

	return value, score, nil
}
