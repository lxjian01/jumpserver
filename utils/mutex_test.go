package utils

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T)  {
	ml := NewMutex()
	ml.TryLock()
	ml.Unlock()
	ml.TryLock()
	ml.Unlock()
	ml.TryLock()
	ml.Unlock()
	fmt.Println("finish")
}

func TestMutex_TryLock(t *testing.T) {

	ml := NewMutex()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		fmt.Println("before r1 Lock, lock is Locked = ", ml.IsLocked())

		ml.Lock()
		defer ml.Unlock()

		fmt.Println("r1 get lock!")

		time.Sleep(5 * time.Second)

		fmt.Println("r1 release lock!")

		wg.Done()
	}()

	go func() {

		fmt.Println("before r2 Lock, lock is Locked = ", ml.IsLocked())

		ml.Lock()
		defer ml.Unlock()

		fmt.Println("r2 get lock!")

		time.Sleep(3 * time.Second)

		fmt.Println("r2 release lock!")

		wg.Done()
	}()

	go func() {
		fmt.Println("before r3 Lock, lock is Locked = ", ml.IsLocked())

		ml.Lock()
		defer ml.Unlock()

		fmt.Println("r3 get lock!")

		time.Sleep(2 * time.Second)

		fmt.Println("r3 release lock!")

		wg.Done()
	}()

	for {
		time.Sleep(1 * time.Second)

		fmt.Println("init routine, lock is Locked = ", ml.IsLocked())
		if ml.TryLock() {
			fmt.Println("yeah, try lock success!")
			ml.Unlock()
			break
		} else {
			fmt.Println("init routine try lock failed!")
		}
	}

	wg.Wait()
	t.Log("test over!")
}

func TestMutex_TryAcquire(t *testing.T) {

	ml := NewMutex()
	go func() {
		for i := 0; i < 1000; i++ {
			go ml.test("111")
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			go ml.test("222")
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			go ml.test("333")
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			go ml.test("444")
		}
	}()
	go func() {
		for {
			//fmt.Println(ml.ops)
			time.Sleep(time.Second)
		}
	}()
	select {}
}

func (m *Mutex) test(name string) {
	if !m.TryAcquire() {
		return
	}
	time.Sleep(time.Second * 10)
	defer func() {
		m.Release()
	}()

	fmt.Println(name)
}
