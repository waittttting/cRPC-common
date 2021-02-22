package clock

import (
	"github.com/go-redis/redis"
	"time"
)

const redisUnLockSingle = 0

type redisLock struct {
	client *redis.Client
	singleChan chan int
	key string
}

func newRedisLock(host string, pwd string, dbIndex int, singleChan chan int) (*redisLock, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
		DB:       dbIndex,
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &redisLock{
		client: client,
		singleChan: singleChan,
	}, nil
}

func (rl *redisLock) Lock(key string, value string, lease int64) (bool, error) {
	// 查看 key 是否存在
	_, err := rl.client.Get(key).Result()
	if err != nil && err.Error() != errRedisNil.Error() {
		return false, err
	}
	// key 已经存在
	if err == nil {
		return false, errKeyExist
	}
	ret, err := rl.client.SetNX(key, value, time.Duration(lease) * time.Second).Result()
	if err != nil {
		return false, err
	}
	if !ret {
		return false, nil
	}
	rl.key = key
	// 续租协程
	go func(singleChan chan int) {
		for {
			select {
			case single := <-singleChan:
				if single == redisUnLockSingle {
					break
				}
			default:
			}
			time.Sleep(time.Duration(lease / 3 * 2) * time.Second)
			ret, err = rl.client.SetNX(key, value, time.Duration(lease) * time.Second).Result()
		}
	}(rl.singleChan)
	return true, nil
}

func (rl *redisLock) UnLock() (bool, error) {
	// 停止续租
	rl.singleChan <- redisUnLockSingle
	_, err := rl.client.Del(rl.key).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rl *redisLock) Close() error {
	return rl.client.Close()
}