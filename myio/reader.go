package myio

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Input struct {
	Num  int
	Data []int
}

func Read(fn string) *Input {
	var (
		num  int
		data []int
	)

	f, err := os.Open(fn)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}
		intVal, err := strconv.Atoi(string(line))
		if err != nil {
			fmt.Println(err)
			break
		}
		if num == 0 {
			num = intVal
		} else {
			data = append(data, intVal)
		}
	}

	if err != nil && err != io.EOF {
		return nil
	}
	return &Input{num, data}
}
