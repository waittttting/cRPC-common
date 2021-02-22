package clock

import (
	"testing"
	"time"
)

func TestEtcdLock(t *testing.T) {

	testKey := "/testLock/key1"
	etcdHosts := []string{"127.0.0.1:2379"}

	lock1, err := EtcdLock(etcdHosts)
	if err != nil {
		t.Errorf("create lock1 error: %v", err)
	}
	lock1.SetLease(5)
	ok, err := lock1.Acquire(testKey)
	if err != nil || !ok {
		t.Errorf("lock1 acquire error: %v", err)
	}
	t.Logf("lock1 locked")
	go func() {
		lock2, err := EtcdLock(etcdHosts)
		if err != nil {
			t.Errorf("create lock2 error: %v", err)
		}
		i := 0
		for {
			ok, err := lock2.Acquire(testKey)
			if err != nil {
				if err.Error() == "key already exist"  {
					t.Logf("lock2 try to get lock, loop count: %v", i)
				} else {
					t.Errorf("lock2 acquire error: %v",  err)
				}
			}
			if ok {
				t.Logf("lock2 get locked, loop count: %v", i)
				break
			}
			time.Sleep(1 * time.Second)
			i ++
		}
	}()
	time.Sleep(10 * time.Second)
	ok, err = lock1.Release()
	if err != nil {
		t.Errorf("lock1 release err: %v", err)
	}
	if ok {
		t.Logf("lock1 release")
	}
	time.Sleep(1 * time.Second)
}
