package twheel

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type TimeWheel struct {
	wheel      []map[interface{}]interface{}
	indexMap   map[interface{}]int
	curIndex   int
	wheelCap   int
	lock       sync.Mutex
	noticeChan chan interface{}
	timeout    int
}

/**
 * @Description: 创建一个新的时间轮
 * @param cap 时间轮中循环数组的大小
 * @param noticeChanLen 元素超时回调队列
 * @return *TimeWheel
 */
func New(cap int, timeout int, noticeChan chan interface{}) (*TimeWheel, error) {

	if cap <= 0 {
		return nil, errors.New("cap can not <= 0")
	}

	if timeout <= 0 {
		return nil, errors.New("timeout can not <= 0")
	}

	return &TimeWheel{
		wheelCap:   cap,
		wheel:      make([]map[interface{}]interface{}, cap),
		curIndex:   0,
		noticeChan: noticeChan,
		indexMap:   make(map[interface{}]int),
		timeout:    timeout,
	}, nil
}

func (tw *TimeWheel) Start() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			tw.curIndex = (tw.curIndex + 1) % tw.wheelCap
			tw.lock.Lock()
			curBucket := tw.wheel[tw.curIndex]
			if len(curBucket) == 0 {
				tw.lock.Unlock()
				continue
			}
			for _, value := range curBucket {
				timer := time.NewTimer(500 * time.Millisecond)
				// 1. 通知
				select {
				case tw.noticeChan <- value:
				case <-timer.C:
					logrus.Errorf("item send to noticeChan timeout")
				}
				// 2. 删除 index map 中的数据
				delete(tw.indexMap, value)
			}
			// 3. 清空 curIndex 上的 map
			tw.wheel[tw.curIndex] = make(map[interface{}]interface{})
			tw.lock.Unlock()
		}
	}()
}

func (tw *TimeWheel) Add(item interface{}) {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	itemIndex := (tw.curIndex + tw.timeout) % tw.wheelCap
	tw.indexMap[item] = itemIndex
	if tw.wheel[itemIndex] == nil {
		tw.wheel[itemIndex] = map[interface{}]interface{}{}
	}
	tw.wheel[itemIndex][item] = item
}

func (tw *TimeWheel) Delete(item interface{}) {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	itemIndex := tw.indexMap[item]
	delete(tw.wheel[itemIndex], item)
	delete(tw.indexMap, item)
}

func (tw *TimeWheel) Refresh(item interface{}) {
	tw.lock.Lock()
	defer tw.lock.Unlock()
	curIndex := tw.indexMap[item]
	delete(tw.wheel[curIndex], item)
	delete(tw.indexMap, item)
	newIndex := (curIndex + tw.timeout) % tw.wheelCap
	tw.indexMap[item] = newIndex
	if tw.wheel[newIndex] == nil {
		tw.wheel[newIndex] = map[interface{}]interface{}{}
	}
	tw.wheel[newIndex][item] = item
}
