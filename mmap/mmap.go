package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	SIZE = 4096
	CL   = byte(10)
)

type CacheBlock []byte

func indexCL(data []byte) int {
	for i := 0; i < len(data); i += 1 {
		if data[i] == CL {
			return i
		}
	}
	return -1
}

func lastIndexCL(data []byte) int {
	for i := len(data) - 1; i > -1; i -= 1 {
		if data[i] == CL {
			return i
		}
	}
	return -1
}

func SnowFoxRead(fn string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		buf  CacheBlock
		line []byte
		d    int
		pre  int
		next int
		left CacheBlock
	)

	file, err := os.Open(fn)
	if err != nil {
		log.Println(err)
		return
	}

	defer file.Close()

	var t int
	fd := file.Fd()
	fi, err := file.Stat()
	fsize := fi.Size()

	dh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	dh.Len = int(SIZE)
	dh.Cap = int(SIZE)
	addr, err := mmap(0, uintptr(fi.Size()), uintptr(PROT_READ), uintptr(MAP_PRIVATE), fd, int64(0))
	addr_bak := addr
	if err != syscall.Errno(0) {
		log.Println(t, err)
		return
	}
	defer func() {
		unmap(addr_bak, uintptr(fi.Size()))
		dh = nil
		buf = nil
	}()

	sf := MakeStringFinder("/v1/product/productView")
	var flag bool
	var cnt int
	var cr_last int
	//var pat_len = len(sf.pattern)
	for fz := fsize; fz > 0; fz -= SIZE {
		dh.Len = int(SIZE)
		dh.Cap = int(SIZE)
		line = nil
		if flag {
			break
		}

		if fz < SIZE {
			dh.Len = int(fz)
			dh.Cap = int(fz)
			flag = true
		}
		dh.Data = addr

		buf = append(left, buf...)
		cr_last = lastIndexCL(buf)
		left = buf[cr_last+1:]
		buf = buf[:cr_last]

		for {
			d = sf.Next(string(buf))
			if d < 0 {
				break
			}
			pre = lastIndexCL(buf[:d])
			next = indexCL(buf[d:])

			if pre == -1 {
				pre = 0
			}

			if next < 0 {
				next = len(buf)
			} else {
				next += d
			}

			line = buf[pre:next]
			_ = line

			cnt += 1
			if len(buf) == next {
				break
			}
			buf = buf[next+1:]
		}

		addr += 4096
		if fz < SIZE {
			break
		}
	}
	fmt.Println(cnt)
}

func main() {
	fmt.Println(time.Now().Unix())

	for i := 23; i < 40; i++ {
		fn := fmt.Sprintf("/2015-08-25_000%d", i)
		//fn := "/tmp/200.lst"
		SnowFoxRead(fn)
	}
	fmt.Println(time.Now().Unix())
}
func mmap(addr, length, prot, flags, fd uintptr, offset int64) (uintptr, error) {
	addr, _, err := syscall.Syscall6(syscall.SYS_MMAP, addr, length, prot, flags, fd, uintptr(offset))
	return addr, err
}

func unmap(addr, l uintptr) error {
	_, _, errno := syscall.Syscall(syscall.SYS_MUNMAP, addr, l, 0)
	if errno != 0 {
		return syscall.Errno(errno)
	}
	return nil
}

type ProtFlags uint

const (
	PROT_NONE  ProtFlags = 0x0
	PROT_READ  ProtFlags = 0x1
	PROT_WRITE ProtFlags = 0x2
	PROT_EXEC  ProtFlags = 0x4
)

type MapFlags uint

const (
	MAP_SHARED    MapFlags = 0x1
	MAP_PRIVATE   MapFlags = 0x2
	MAP_FIXED     MapFlags = 0x10
	MAP_ANONYMOUS MapFlags = 0x20
	MAP_GROWSDOWN MapFlags = 0x100
	MAP_LOCKED    MapFlags = 0x2000
	MAP_NONBLOCK  MapFlags = 0x10000
	MAP_NORESERVE MapFlags = 0x4000
	MAP_POPULATE  MapFlags = 0x8000
)

type SyncFlags uint

const (
	MS_SYNC       SyncFlags = 0x4
	MS_ASYNC      SyncFlags = 0x1
	MS_INVALIDATE SyncFlags = 0x2
)

type AdviseFlags uint

const (
	MADV_NORMAL     AdviseFlags = 0x0
	MADV_RANDOM     AdviseFlags = 0x1
	MADV_SEQUENTIAL AdviseFlags = 0x2
	MADV_WILLNEED   AdviseFlags = 0x3
	MADV_DONTNEED   AdviseFlags = 0x4
	MADV_REMOVE     AdviseFlags = 0x9
	MADV_DONTFORK   AdviseFlags = 0xa
	MADV_DOFORK     AdviseFlags = 0xb
)

type StringFinder struct {
	pattern        string
	badCharSkip    [256]int
	goodSuffixSkip []int
}

func MakeStringFinder(pattern string) *StringFinder {
	f := &StringFinder{
		pattern:        pattern,
		goodSuffixSkip: make([]int, len(pattern)),
	}
	last := len(pattern) - 1

	for i := range f.badCharSkip {
		f.badCharSkip[i] = len(pattern)
	}

	for i := 0; i < last; i++ {
		f.badCharSkip[pattern[i]] = last - i
	}

	lastPrefix := last
	for i := last; i >= 0; i-- {
		if strings.HasPrefix(pattern, pattern[i+1:]) {
			lastPrefix = i + 1
		}
		f.goodSuffixSkip[i] = lastPrefix + last - i
	}
	for i := 0; i < last; i++ {
		lenSuffix := longestCommonSuffix(pattern, pattern[1:i+1])
		if pattern[i-lenSuffix] != pattern[last-lenSuffix] {
			f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
		}
	}

	return f
}

func longestCommonSuffix(a, b string) (i int) {
	for ; i < len(a) && i < len(b); i++ {
		if a[len(a)-1-i] != b[len(b)-1-i] {
			break
		}
	}
	return
}

func (f *StringFinder) Next(text string) int {
	i := len(f.pattern) - 1
	for i < len(text) {
		j := len(f.pattern) - 1
		for j >= 0 && text[i] == f.pattern[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1 // match
		}
		i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (f *StringFinder) Grep(src string) bool {
	match := f.Next(src)
	if match == -1 {
		return false
	}
	return true
}
