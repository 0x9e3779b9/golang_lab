package main

import "fmt"
import "unsafe"
import "reflect"

func main() {
	s := []int{1, 2, 3, 4, 5, 6}
	p := unsafe.Pointer(&s)
	size := *(*int)(unsafe.Pointer(uintptr(p) + uintptr(unsafe.Alignof(reflect.SliceHeader.Len))))
	fmt.Printf("sizeof = %d\n", size)
	//p2 := (*[]byte)(unsafe.Pointer(uintptr(p)))
	//fmt.Printf("str = %v\n", string(*p2))
}
