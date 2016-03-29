package main

import "fmt"
import "unsafe"

func main() {
	var str = "ABC12"
	p := unsafe.Pointer(&str)
	size := *(*int)(unsafe.Pointer(uintptr(p) + uintptr(8)))
	fmt.Printf("sizeof = %d\n", size)
	p2 := (*[]byte)(unsafe.Pointer(uintptr(p)))
	fmt.Printf("str = %v\n", string(*p2))
}
