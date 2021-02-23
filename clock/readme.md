# 分布式锁

## 基于 etcd 的 CP 分布式锁

```go
EtcdLock
```
调用方法
```go
lock1, err := EtcdLock(etcdHosts)

// 获取非阻塞的锁
ok, err := lock1.Acquire(testKey1)

// 获取阻塞一段时长的锁
ok, err = lock1.AcquireWithRetry(testKey1, interval, maxRetry)
```

## 基于 redis 的 AP 分布式锁

```go
RedisLock
```
调用方法
```go
lock1, err := RedisLock(redisHost, "", 0)

// 获取非阻塞的锁
ok, err := lock1.Acquire(testKey1)

// 获取阻塞一段时长的锁
ok, err = lock1.AcquireWithRetry(testKey1, interval, maxRetry)
```

