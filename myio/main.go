package main

import (
	"strconv"
	"strings"
)

func getIntVal(str string) uint8 {
	ret, err := strconv.Atoi(str)
	if err != nil {
		return uint8(0)
	}
	return uint8(ret)
}

func getVal(str string) (ret uint64) {
	list := strings.Split(str, ".")
	ret |= uint64(getIntVal(list[0])) << 24
	println(ret)
	ret |= uint64(getIntVal(list[1])) << 16
	println(ret)
	ret |= uint64(getIntVal(list[2])) << 8
	println(ret)
	ret |= uint64(getIntVal(list[3]))
	return
}

func main() {
	println(getVal("222.11.12.13"))
}
