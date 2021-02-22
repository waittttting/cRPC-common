package clock

type RedisLock struct {}

func (al *RedisLock) Lock(key string, value string, lease int64) (bool, error) {
	return true, nil
}

func (al *RedisLock) UnLock(key string) (bool, error) {
	return true, nil
}
