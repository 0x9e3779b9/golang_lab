package notify

import (
	"fmt"
	"log"

	"github.com/howeyc/fsnotify"
	"golang_lab/mytask/message"
)

const (
	EOF = "EOF"
)

type LogWatcher struct {
	watcher   *fsnotify.Watcher
	filePath  string
	msgChPool chan []byte
	eofCh     chan struct{}
	out       chan *message.Message
}

func NewLogWatcher(fn string, ch chan *message.Message) (w *LogWatcher, err error) {
	wt, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	w = &LogWatcher{
		watcher:   wt,
		filePath:  fn,
		msgChPool: make(chan []byte, 1024),
		eofCh:     make(chan struct{}),
		out:       ch,
	}
	w.watcher.Watch(fn)
	return
}

func (w *LogWatcher) Monitor() {
	size := 0
	for {
		select {
		case ev := <-w.watcher.Event:
			if ev.IsModify() {
				sizeN := getFileSize(w.filePath)
				if sizeN > size {
					w.msgChPool <- readBytes(w.filePath, size, sizeN-1)
				}
				size = sizeN
			}
		case err := <-w.watcher.Error:
			log.Println("error:", err)
		}
	}
}

func (w *LogWatcher) Accept() {
	var lines []byte
	for {
		select {
		case msg := <-w.msgChPool:
			lines = append(lines, msg...)
			index := 0
			for index > -1 {
				index = getBytesIndex(lines, byte(10))
				if index > -1 {
					w.out <- &message.Message{
						w.filePath,
						string(lines[:index]),
					}
					if string(lines[:index]) == EOF {
						fmt.Println("OVER")
						w.eofCh <- struct{}{}
					}
					lines = lines[index+1:]
				}
			}
		}
	}
}
