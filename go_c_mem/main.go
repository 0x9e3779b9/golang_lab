package main

/*
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

char* newData(size_t len)
{
	char *ptr = (char *)malloc(sizeof(char)*len);
	memcpy(ptr,"hello",6);
	return ptr;
}
*/
import "C"

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

func main() {
	ptr := C.newData(1024)
	p := unsafe.Pointer(ptr)
	var data []byte
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	slice.Data = uintptr(p)
	slice.Cap = 1024
	slice.Len = 1024

	// 对于使用c申请的内存,将其加入到GC中
	runtime.SetFinalizer(slice, func(s *reflect.SliceHeader) {
		println("GC called")
		C.free(unsafe.Pointer(s.Data))
	})

	fmt.Printf("%s\n", string(data))
	time.Sleep(time.Second)
	runtime.GC()
	fmt.Printf("%s\n", string(data))
	runtime.GC()
	time.Sleep(time.Second * 2)
}
