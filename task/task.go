package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type A struct {
	cnt     int
	wg      *sync.WaitGroup
	pauseCh chan struct{}
	goCh    chan struct{}
	stopCh  chan struct{}
}

func (a *A) do() {
	index := 0
	for index < 4 {
		time.Sleep(time.Second * 1)
		print("do")
		println(index)
		index++
	}
}

func (a *A) Run() {
	var tag int32
	var num int
	a.wg.Add(1)
	defer a.wg.Done()
	for {
		select {
		case <-a.pauseCh:
			println("pause")
			atomic.StoreInt32(&tag, 1)
		case <-a.goCh:
			println("go")
			atomic.StoreInt32(&tag, 0)
		case <-a.stopCh:
			println("stop")
			break
		default:
			if atomic.LoadInt32(&tag) == 0 {
				println("default")
				a.do()
				num++
				print("finished")
				println(num)
				if num >= a.cnt {
					break
				}
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}
}

func (a *A) Pause() {
	a.pauseCh <- struct{}{}
}

func (a *A) Go() {
	a.goCh <- struct{}{}
}

func (a *A) Stop() {
	a.stopCh <- struct{}{}
}

func (a *A) Wait() {
	a.wg.Wait()
}

func main() {
	a := &A{
		goCh:    make(chan struct{}),
		pauseCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
		wg:      new(sync.WaitGroup),
		cnt:     4,
	}

	go a.Run()
	time.Sleep(time.Second * 2)
	a.Pause()
	time.Sleep(time.Second * 2)
	a.Go()
	a.Wait()
}
