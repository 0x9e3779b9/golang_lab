/*
 * Copyright (C)  Dongchen Yang
 * Copyright (C) xiaomi.com
 */
package main

import (
	"encoding/binary"
	"os"
	"sync"
)

type Partition struct {
	f      *os.File
	offset uint64
	lock   *sync.Mutex
}

func newPartition(fn string) (p *Partition, err error) {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return
	}
	p = &Partition{
		f:    f,
		lock: new(sync.Mutex),
	}
	return
}

func (p *Partition) push(data []byte) {
	p.lock.Lock()
	n, _ := p.f.Write(data)
	p.offset += uint64(n)
	p.lock.Unlock()
}

func (p *Partition) pull(offset int64) {
	ret := make([]byte, 9)
	n, err := p.f.ReadAt(ret, offset)
	if err != nil || n != 9 {
		return
	}
	all := binary.BigEndian.Uint32(ret[:4])
	all -= 9
	ret = make([]byte, int(all))
	p.f.ReadAt(ret, offset+9)
}

var tmpl []byte

func main() {
	size := 1024
	test := make([]byte, size)
	for i := 0; i < size; i++ {
		test[i] = byte('a')
	}
	tmpl = make([]byte, size+9)
	binary.BigEndian.PutUint32(tmpl[:4], uint32(size+9))
	tmpl[4] = byte('A')
	binary.BigEndian.PutUint32(tmpl[5:9], uint32(0x1024))
	copy(tmpl[9:], test)
	st, err := newPartition("file_test")
	if err != nil {
		println(err.Error())
		return
	}
	for i := 0; i < 100000; i++ {
		st.push(tmpl)
	}
	st.pull(0)
	st.f.Close()
}
