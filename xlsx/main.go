package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

func Test() {
	xlFile := xlsx.NewFile()
	sh := xlFile.AddSheet("test")
	header := sh.AddRow()
	header.Cells = []*xlsx.Cell{
		&xlsx.Cell{Value: "Name"},
		&xlsx.Cell{Value: "ID"},
		&xlsx.Cell{Value: "Cost"},
		&xlsx.Cell{Value: "Result"},
	}
	for i := 0; i < 10; i++ {
		row := sh.AddRow()
		row.Cells = []*xlsx.Cell{
			&xlsx.Cell{Value: "Mike"},
			&xlsx.Cell{Value: fmt.Sprintf("%d", i+1)},
			&xlsx.Cell{Value: fmt.Sprintf("%d", i+20)},
			&xlsx.Cell{Value: "pass"},
		}
	}
	xlFile.Save("/Users/dongchen/dumps/demo.xlsx")
}

func main() {
	Test()
}
