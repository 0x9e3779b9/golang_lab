package main

import (
	"fmt"
	"log"
	"os"

	"github.com/howeyc/fsnotify"
)

func getFileSize(fn string) int {
	fileInfo, err := os.Stat(fn)
	if err != nil {
		log.Fatal(err)
	}
	return int(fileInfo.Size())
}

func readBytes(fname string, start, end int) []byte {
	fh, _ := os.Open(fname)
	defer fh.Close()

	fh.Seek(int64(start), 0)
	size := end - start + 1
	buff := make([]byte, size)
	fh.Read(buff)
	return buff
}

func getBytesIndex(bts []byte, sep byte) int {
	for index := range bts {
		if sep == bts[index] {
			return index
		}
	}
	return -1
}

type LogWatcher struct {
	watcher   *fsnotify.Watcher
	msgChPool map[string]chan []byte
}

func NewLogWatcher() (w *LogWatcher, err error) {
	wt, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	w = &LogWatcher{
		watcher:   wt,
		msgChPool: make(map[string]chan []byte),
	}
	return
}

func (w *LogWatcher) AddPath(fn string) {
	msgChPool[fn] = make(chan []byte, 1024)
}

func (w *LogWatcher) RemovePath(fn string) {
	w.watcher.RemoveWatch(fn)
}

func (w *LogWatcher) Monitor() {
}

func (w *LogWatcher) Broadcast() {

}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	fname := "/tmp/fsnotify_test.log"
	done := make(chan bool)
	msgCh := make(chan []byte, 1024)

	go func() {
		var lines []byte
		for {
			select {
			case msg := <-msgCh:
				lines = append(lines, msg...)
				index := 0
				for index > -1 {
					index = getBytesIndex(lines, byte(10))
					if index > -1 {
						fmt.Println(string(lines[:index]))
						lines = lines[index+1:]
					}
				}
			}
		}
	}()

	go func() {
		size := getFileSize(fname)
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() {
					sizeN := getFileSize(fname)
					if sizeN > size {
						msgCh <- readBytes(fname, size, sizeN-1)
					}
					size = sizeN
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(fname)
	if err != nil {
		log.Fatal(err)
	}

	// Hang so program doesn't exit
	<-done

	/* ... do stuff ... */
	watcher.Close()
}
