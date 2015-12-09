package notify

import (
	"fmt"

	"golang_lab/mytask/message"
)

type FSnotifyManager struct {
	path  chan string
	outCh chan *message.Message
}

func NewFSnotifyManager(size int, ch chan *message.Message) *FSnotifyManager {
	return &FSnotifyManager{
		path:  make(chan string, size),
		outCh: ch,
	}
}

func (f *FSnotifyManager) Push(fn string) {
	f.path <- fn
}

func (f *FSnotifyManager) Run() {
	for {
		select {
		case fn := <-f.path:
			go f.do(fn)
		}
	}
}

func (f *FSnotifyManager) do(fn string) {
	w, err := NewLogWatcher(fn, f.outCh)
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
