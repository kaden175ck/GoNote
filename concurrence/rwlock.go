package concurrence

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var lock sync.Mutex //sync.Mutex相当于是写锁

func inc3() {
	lock.Lock()   //加写锁
	n++           //任一时刻，只有一个协程能进入临界区域
	lock.Unlock() //释放写锁
}

func Lock() {
	const P = 1000 //开大量协程才能把脏写问题测出来
	wg := sync.WaitGroup{}
	wg.Add(P)
	for i := 0; i < P; i++ {
		go func() {
			defer wg.Done()
			inc3()
		}()
	}
	wg.Wait()
	fmt.Printf("finally n=%d\n", n)
	fmt.Println("===========================")
}

var (
	mu sync.RWMutex
)

// 读锁是可重入的，即同一个协程可以多次获得同一把读锁（之前获得的还没释放）
func ReentranceRLock(n int) {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println(n)
	if n > 0 {
		ReentranceRLock(n - 1)
	}
	time.Sleep(1 * time.Second)
}

// 写锁是不可重入的
func ReentranceWLock(n int) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(n)
	if n > 0 {
		ReentranceWLock(n - 1)
	}
	time.Sleep(1 * time.Second)
}

func RLockExclusion() {
	mu.RLock() //获得锁写，其他协程不可以获得读锁，也不可以获得写锁
	go func() {
		mu.Lock() // 如果有写锁在等待，那么读锁也要阻塞
		defer mu.Unlock()
		fmt.Println("子协程也获得了写锁")
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mu.RLock() // 如果有写锁在等待，那么读锁也要阻塞
		defer mu.RUnlock()
		fmt.Println("子协程也获得了读锁")
	}()

	time.Sleep(1 * time.Second)
	mu.RUnlock()
	time.Sleep(1 * time.Second)
}

func WLockExclusion() {
	mu.Lock() //获得锁写，其他协程不可以获得读锁，也不可以获得写锁
	defer mu.Unlock()
	go func() {
		mu.RLock()
		defer mu.RUnlock()
		fmt.Println("子协程也获得了读锁")
	}()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("子协程也获得了写锁")
	}()
	time.Sleep(1 * time.Second)
}

// 排队等待锁
func LockQueue() {
	mu.RLock()

	go func() {
		mu.RLock()
		defer mu.RUnlock()
		log.Println("协程1获得了读锁")
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mu.RLock()
		defer mu.RUnlock()
		log.Println("协程2获得了读锁")
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		log.Println("协程3获得了写锁")
	}()

	go func() {
		time.Sleep(300 * time.Millisecond)
		mu.Lock()
		defer mu.Unlock()
		log.Println("协程4获得了写锁")
	}()

	//按照【排队】原则，协程1、2、3、4依次获得锁

	time.Sleep(1 * time.Second)
	mu.RUnlock()
	time.Sleep(1 * time.Second)
}

func ReadWriteRace() {
	mu.RLock()
	log.Println("协程1获得了读锁")
	go func() {
		mu.RLock()
		defer mu.RUnlock()
		log.Println("协程2获得了读锁")
		time.Sleep(1 * time.Second)
	}()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		log.Println("协程3获得了写锁")
		time.Sleep(1 * time.Second)
	}()
	time.Sleep(1 * time.Second)
	mu.RUnlock()
	log.Println("协程1释放了读锁")
	time.Sleep(2 * time.Second)
}

/**
总共只会出现以下2种情况：
2025/07/23 20:45:14 协程1获得了读锁
2025/07/23 20:45:14 协程2获得了读锁
2025/07/23 20:45:15 协程1释放了读锁
2025/07/23 20:45:15 协程3获得了写锁
协程2先于协程3执行，所以协程2顺利获得了读锁，不需要等协程1释放读锁。


2025/07/23 20:45:08 协程1获得了读锁
2025/07/23 20:45:09 协程1释放了读锁
2025/07/23 20:45:09 协程3获得了写锁
2025/07/23 20:45:10 协程2获得了读锁
协程3先于协程2执行，协程1释放读锁后，协程3会先获得锁（【排队】机制，协程3比协程2执行），既然协程3获得了写锁，那么协程2就获得不到读锁。
*/
