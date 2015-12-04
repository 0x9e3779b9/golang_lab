package main

func main() {
	Init(2)
	go Run()
	Push("/tmp/a")
	Push("/tmp/b")
	select {}
}
