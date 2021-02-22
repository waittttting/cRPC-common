package clock

import (
	"context"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type etcdLock struct {

	client *clientv3.Client
	leaseID clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
}

func newEtcdLock(hosts []string) (*etcdLock, error) {

	config := clientv3.Config{
		Endpoints:   hosts,
		DialTimeout: 5 * time.Second,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &etcdLock{
		client: client,
	}, nil
}


func (el *etcdLock) Lock(key string, value string, lease int64) (bool, error) {

	if key == "" {
		return false, errors.New("key len is 0")
	}
	// 查看 key 是否存在
	gr, err := el.client.Get(context.Background(), key)
	if err != nil {
		return false, err
	}
	if len(gr.Kvs) > 0 {
		return false, errors.New("key already exist")
	}
	// 申请租约
	resp, err := el.client.Grant(context.Background(), lease)
	if err != nil {
		return false, err
	}
	// 使用租约注册 key
	_, err = el.client.Put(context.Background(), key, value, clientv3.WithLease(resp.ID))
	if err != nil {
		return false, err
	}
	// 设置续租 定期发送需求请求
	leaseRespChan, err := el.client.KeepAlive(context.Background(), resp.ID)
	if err != nil {
		return false, err
	}
	el.leaseID = resp.ID
	el.keepAliveChan = leaseRespChan
	return true, nil
}

func (el *etcdLock) UnLock() (bool, error) {
	if _, err := el.client.Revoke(context.Background(), el.leaseID); err != nil {
		return false, err
	}
	el.client.Close()
	return true, nil
}
