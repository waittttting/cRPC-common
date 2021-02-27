package twheel

import (
	"testing"
	"time"
)

var timeout = 5

func TestNew(t *testing.T) {

	ch := make(chan interface{}, 20)
	go func() {
		i := 0
		for msg := range ch {
			i++
			t.Logf("i = %d, %v", i, msg)
		}
	}()

	_, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
}

func testNew(ch chan interface{}) (*TimeWheel, error) {

	return New(20, timeout, ch)
}

type testIns struct {
}

func TestTimeWheel_Add(t *testing.T) {

	var addTime time.Time
	ch := make(chan interface{}, 20)
	go func() {
		i := 0
		for msg := range ch {
			i++
			t.Logf("i = %d, %v", i, msg)
			diff := int(time.Now().Sub(addTime).Seconds())
			if diff != timeout {
				t.Error("timeout was wrong")
			}
		}
	}()

	tw, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
	tw.Start()
	tw.Add(new(testIns))
	addTime = time.Now()
	time.Sleep(10 * time.Second)
}

func TestTimeWheel_Refresh(t *testing.T) {

	var addTime time.Time
	ch := make(chan interface{}, 20)
	go func() {
		i := 0
		for msg := range ch {
			i++
			t.Logf("i = %d, %v", i, msg)
			diff := int(time.Now().Sub(addTime).Seconds())
			if diff != timeout*2 {
				t.Error("timeout was wrong")
			}
		}
	}()

	tw, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
	tw.Start()
	ins := new(testIns)
	tw.Add(ins)
	addTime = time.Now()
	tw.Refresh(ins)
	time.Sleep(15 * time.Second)
}

func TestTimeWheel_Delete(t *testing.T) {

	ch := make(chan interface{}, 20)
	go func() {
		for range ch {
			t.Error("timeout was wrong")
		}
	}()

	tw, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
	tw.Start()
	ins := new(testIns)
	tw.Add(ins)
	tw.Delete(ins)
	time.Sleep(10 * time.Second)
}
