/*********************************************
*  Author:  Dongchen Yang
*  Last modified: 2015-12-04 11:47
*  Filename: main.go
*  Description:
*********************************************/
package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type proxy struct {
	Host  string
	Port  string
	Local string
}

func runProxy(list []proxy) {
	wg := sync.WaitGroup{}
	for _, v := range list {
		wg.Add(1)
		go func() {
			eachServer(v.Host, v.Port, v.Local)
			wg.Done()
		}()
	}
	wg.Wait()
}

func eachServer(host string, port string, local string) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		fmt.Print("listen tcp :%s", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		tc, err := l.Accept()
		if err != nil {
			fmt.Printf("accept tcp conn :%s", err.Error())
			tc.Close()
			continue
		}

		go eachConn(local, tc)
	}
}

func eachConn(local string, tc net.Conn) {
	uc, err := net.Dial("unix", local)
	if err != nil {
		fmt.Printf("get unix conn :%s", err.Error())
		uc.Close()
		return
	}
	go io.Copy(tc, uc)
	go io.Copy(uc, tc)
}

func main() {
	list := []proxy{proxy{Host: "192.168.8.101", Port: "3306", Local: "/home/work/data/tmp/mysql3306.sock"}}
	runProxy(list)
}
