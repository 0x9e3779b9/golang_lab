package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func dup2(fn string, buf *bufio.Reader) {
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(line))
		f.WriteString(string(line) + "\n")
	}
}

func stdOut(buf *bufio.Reader) {
	dup2("test.log", buf)
}

func stdErr(buf *bufio.Reader) {
	dup2("test.err", buf)
}

func test() {
	cmd := exec.Command("/usr/local/bin/python", "test.py")
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	outBuf := bufio.NewReader(outPipe)
	errBuf := bufio.NewReader(errPipe)
	go stdOut(outBuf)
	go stdErr(errBuf)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cmd.Process.Pid)
	defer cmd.Process.Kill()
	defer cmd.Process.Release()
	cmd.Wait()

}

func main() {
	go test()
	time.Sleep(time.Second * 60)
}
