package main

import (
	"fmt"
	"golang_lab/myio"
	"time"
)

func QuickSort(src []int, start, end int) {
	if end <= start {
		return
	}
	key := src[(start+end)/2]
	var i, j = start, end
	for i < j {
		for src[i] < key && i < end {
			i++
		}
		for src[j] > key && j > start {
			j--
		}

		if i < j {
			src[i], src[j] = src[j], src[i]
		}
	}
	key, src[j] = src[j], key
	QuickSort(src, start, j-1)
	QuickSort(src, j+1, end)
}

func InsertSort(src []int) {
	for i := 0; i < len(src)-1; i++ {
		for j := i + 1; j > 0; j-- {
			if src[j] < src[j-1] {
				src[j], src[j-1] = src[j-1], src[j]
			}
		}
	}
}

func Prime(max uint64) {
	flag := make([]bool, max)
	var i uint64
	for i = 3; i < max; i += 2 {
		flag[i] = true
	}
	flag[2] = true
	start := time.Now()
	for i = 3; i < max; i += 2 {
		if flag[i] {
			var j uint64 = i * i
			var n uint64 = i
			for j <= max {
				flag[j] = false
				for max/j >= i {
					flag[j*i] = false
					j = j * i
				}
				n += 2
				for !flag[n] {
					n += 2
				}
				j = i * n
			}
		}
	}
	fmt.Println(time.Now().Sub(start))

	cnt := 0
	for index := range flag {
		if flag[index] {
			cnt++
		}
	}
	fmt.Printf("[0,%d] has %d prime num\n", max, cnt)
}

func main() {
	input := myio.Read("sample.lst")
	if input == nil {
		fmt.Println("ERRR")
		return
	}
	//QuickSort(src, 0, len(src)-1)
	InsertSort(input.Data)
	fmt.Println(input.Data)
	Prime(100000000)
	myio.FlushToFile("result.lst", input.Data)
}
