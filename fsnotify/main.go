package main

import (
	"fmt"
)

func main() {
	w, err := NewLogWatcher("/tmp/fsnotify_test.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	go w.Accept()
	go w.Monitor()
	<-w.eofCh
}

// func test() {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fname := "/tmp/fsnotify_test.log"
// 	done := make(chan bool)
// 	msgCh := make(chan []byte, 1024)

// 	go func() {
// 		var lines []byte
// 		for {
// 			select {
// 			case msg := <-msgCh:
// 				lines = append(lines, msg...)
// 				index := 0
// 				for index > -1 {
// 					index = getBytesIndex(lines, byte(10))
// 					if index > -1 {
// 						fmt.Println(string(lines[:index]))
// 						lines = lines[index+1:]
// 					}
// 				}
// 			}
// 		}
// 	}()

// 	go func() {
// 		size := getFileSize(fname)
// 		for {
// 			select {
// 			case ev := <-watcher.Event:
// 				if ev.IsModify() {
// 					sizeN := getFileSize(fname)
// 					if sizeN > size {
// 						msgCh <- readBytes(fname, size, sizeN-1)
// 					}
// 					size = sizeN
// 				}
// 			case err := <-watcher.Error:
// 				log.Println("error:", err)
// 			}
// 		}
// 	}()

// 	err = watcher.Watch(fname)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Hang so program doesn't exit
// 	<-done

// 	/* ... do stuff ... */
// 	watcher.Close()
// }
