package clock

import (
	"testing"
	"time"
)

var redisHost = "localhost:6379"

func TestCreateRedisLock(t *testing.T) {

	lock1, err := RedisLock(redisHost, "", 0)
	if err != nil {
		t.Fatalf("lock1 create err %v", err)
	}
	lock1.SetLease(3)
	ok, err := lock1.Acquire(testKey1)
	if err != nil || !ok {
		t.Fatalf("lock1 acquire error: %v", err)
	}
	if ok {
		t.Log("lock1 locked")
	}
	time.Sleep(10 * time.Second)
	ok, err = lock1.Release()
	if err != nil {
		t.Fatalf("lock1 release error: %v", err)
	}
	if ok {
		t.Log("lock1 release")
	}
}