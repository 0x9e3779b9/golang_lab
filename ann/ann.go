package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	"golang_lab/ann/core"
)

const (
	ZEROI64 = int64(0)
	ZEROF64 = float64(0)
)

type AnnTrainer struct {
	Hidden   []int64
	Model    []*core.Matrix
	MaxLable int64
}

func (a *AnnTrainer) LoadModel(path string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	byt, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	lines := strings.Split(string(byt), "\n")
	layerNum, err := strconv.Atoi(lines[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	a.Model = make([]*core.Matrix, layerNum)
	lines = lines[1:]

	var index int
	for index < layerNum {
		row, err := strconv.ParseInt(lines[0], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		col, err := strconv.ParseInt(lines[1], 10, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("has %d * %d matrix\n", row, col)
		lines = lines[2:]
		a.Model[index] = core.NewMatrix()
		for i := int64(0); i < row; i++ {
			a.Model[index].Data[i] = core.NewVector()
			for j := int64(0); j < col; j++ {
				a.Model[index].SetValue(j, i, parseFloat64(lines[j]))
			}
			lines = lines[col:]
		}
		index++
	}
}

func (a *AnnTrainer) Predict(sample *core.Sample) {
	x := core.NewVector()
	y := core.NewVector()
	//z := core.NewVector()
	for i := ZEROI64; i < a.Hidden[0]; i++ {
		sum := ZEROF64
		for _, f := range sample.Features {
			sum += f.Value * a.Model[0].Data[i].GetValue(f.Id)
		}
		x.Data[i] = sigMoid(sum)
	}
	x.Data[a.Hidden[0]] = 1
	for i := ZEROI64; i < a.Hidden[1]; i++ {
		sum := ZEROF64
		for j := ZEROI64; j <= a.Hidden[0]; j++ {
			sum += x.GetValue(j) * a.Model[1].GetValue(j, i)
		}
		y.Data[i] = sigMoid(sum)
	}
	y.Data[a.Hidden[1]] = 1

	sum := ZEROF64
	for j := ZEROI64; j <= a.Hidden[1]; j++ {
		sum += y.GetValue(j) * a.Model[2].GetValue(j, 0)
	}
	sample.Prediction = sum
}

func parseFloat64(str string) (ret float64) {
	var err error
	ret, err = strconv.ParseFloat(str, 10)
	if err != nil {
		println(err)
		return
	}
	return
}

func sigMoid(d float64) float64 {
	return 1 / (1 + math.Exp(-1*d))
}

func main() {
	a := AnnTrainer{
		Hidden: []int64{200, 60},
	}
	a.LoadModel("ann_dump")
	for index := range a.Model {
		fmt.Printf("L%d ----- start ------\n", index+1)
		for r := range a.Model[index].Data {
			//fmt.Printf("Row%d has %v\n", r, a.Model[index].Data[r].Data)
			_ = r
		}
	}

}
