package main

import (
	"log"
	"os"
)

func getFileSize(fn string) int {
	fileInfo, err := os.Stat(fn)
	if err != nil {
		log.Fatal(err)
	}
	return int(fileInfo.Size())
}

func readBytes(fname string, start, end int) []byte {
	fh, _ := os.Open(fname)
	defer fh.Close()

	fh.Seek(int64(start), 0)
	size := end - start + 1
	buff := make([]byte, size)
	fh.Read(buff)
	return buff
}

func getBytesIndex(bts []byte, sep byte) int {
	for index := range bts {
		if sep == bts[index] {
			return index
		}
	}
	return -1
}
