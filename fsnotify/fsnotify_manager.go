package main

import (
	"fmt"
)

type FsnotifyManager struct {
	path chan string
}

var defaultMgr *FsnotifyManager

func Init(size int) {
	defaultMgr = &FsnotifyManager{
		path: make(chan string, size),
	}
}

func Push(fn string) {
	defaultMgr.path <- fn
}

func Run() {
	for {
		select {
		case fn := <-defaultMgr.path:
			go do(fn)
		}
	}
}

func do(fn string) {
	w, err := NewLogWatcher(fn)
	if err != nil {
		fmt.Println(err)
		return
	}
	go w.Accept()
	go w.Monitor()
	fmt.Printf("wait for %s\n", fn)
	<-w.eofCh
	fmt.Printf("finish %s\n", fn)
}
