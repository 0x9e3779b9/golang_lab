/*
 * Copyright (C)  Eric
 * Copyright (C) dup2snow@gmail.com
 */
package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func sendFile(conn *net.TCPConn, fn string, offset int64) {
	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	nf, err := conn.File()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	dst := nf.Fd()
	for i := 0; i < 1000000; i++ {
		_, err = syscall.Sendfile(int(dst), int(f.Fd()), &offset, 1033)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
	}
	conn.CloseWrite()
}

func serve() {
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:10011")
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go sendFile(conn, "file_test", 0)
	}
}

func client() {
	laddr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:0")
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:10011")
	conn, err := net.DialTCP("tcp4", laddr, addr)
	if err != nil {
		println(err.Error())
		return
	}
	var d []byte = make([]byte, 4096)
	println(time.Now().Unix())
	for {
		_, err := conn.Read(d)
		if err != nil {
			break
		}
	}
	println(time.Now().Unix())
}

func main() {
	go serve()
	time.Sleep(time.Second * 1)
	for i := 0; i < 50; i++ {
		go client()
	}
	time.Sleep(time.Second * 100)
}
