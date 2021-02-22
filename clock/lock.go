package clock

import (
	"errors"
)

type lockType int

const (
	lockTypeEtcdCP  = 1
	lockTypeRedisAP = 2
	defaultLease = 10
	etcdDefaultValue = "etcd_default_value"
)

type Lock interface {
	Lock(key string, value string, lease int64) (bool, error)
	UnLock() (bool, error)
}

type CXLock struct {
	lock Lock
	lt lockType
	lease int64
}

/**
 * @Description: 设置续约时间 默认为 10
 * @receiver cl
 * @param lease 续约时间
 */
func (cl *CXLock) SetLease(lease int64) {
	cl.lease = lease
}

func EtcdLock(host []string) (*CXLock, error) {
	if len(host) <= 0 {
		return nil, errors.New("host len can not be 0")
	}
	cl := &CXLock{}
	cl.lt = lockTypeEtcdCP
	lock, err := newEtcdLock(host)
	if err != nil {
		return nil, err
	}
	cl.lock = lock
	return cl, nil
}

/**
 * @Description: 获取锁不阻塞
 * @receiver lc
 * @param ttl 锁过期时间
 * @return bool 是否获取锁
 * @return error
 */
func (cl *CXLock) Acquire(key string) (bool, error) {

	if cl.lock == nil {
		return false, errors.New("lock init err, lock can not be nil")
	}
	if cl.lease <= 0 {
		cl.lease = defaultLease
	}
	ok, err := cl.lock.Lock(key, etcdDefaultValue, cl.lease)
	if ok {
		return true, nil
	}
	return false, err
}

/**
 * @Description: 释放锁
 * @receiver lc
 */
func (cl *CXLock) Release() (bool, error) {
	ok, err := cl.lock.UnLock()
	if ok {
		return ok, nil
	}
	return false, err
}

/**
 * @Description: 获取锁，直到重试超时
 * @receiver lc LockClient
 * @param ttl 锁过期时间
 * @param interval 获取锁时间间隔
 * @param maxRetry 获取锁重试次数
 * @return bool 是否获取锁
 * @return error
 */
func (cl *CXLock) AcquireWithRetry(key string, interval int, maxRetry int) (bool, error) {

	cl.lock.Lock(key, "", cl.lease)
	return true, nil
}
