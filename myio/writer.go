package myio

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func FlushToFile(fn string, dst []int) (err error) {
	var buf bytes.Buffer
	for _, d := range dst {
		buf.WriteString(fmt.Sprintf("%d\n", d))
	}
	if err = ioutil.WriteFile(fn, buf.Bytes(), 0644); err != nil {
		return
	}
	return
}
