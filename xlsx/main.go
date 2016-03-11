package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/tealeg/xlsx"
)

var (
	green = &xlsx.Style{Fill: *xlsx.NewFill("solid", "0000FF00", "00FF0000"), ApplyFill: true}
	red   = &xlsx.Style{Fill: *xlsx.NewFill("solid", "00FF0000", "FF000000"), ApplyFill: true}
)

func Test() {
	xlFile := xlsx.NewFile()
	sh := xlFile.AddSheet("result")
	sh.SetColWidth(1, 2, 50)
	sh.MaxCol = 2
	header := sh.AddRow()
	rCell := &xlsx.Cell{Value: "结果"}
	rCell.SetStyle(red)
	header.Cells = []*xlsx.Cell{
		&xlsx.Cell{Value: "编号"},
		rCell,
	}
	for i := 0; i < 10; i++ {
		row := sh.AddRow()
		rs := &xlsx.Cell{Value: "这是一个悲伤的故事，一旦这样就那样。。。。。。\n好的没问题！"}
		rs.SetStyle(green)
		row.Cells = []*xlsx.Cell{
			&xlsx.Cell{Value: "aaaaaa"},
			rs,
		}
	}
	xlFile.Save("aaa.xlsx")
}

func readFileToBytes() {
	b, err := ioutil.ReadFile("~/dumps/demo.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(base64.StdEncoding.EncodeToString(b))
}

func main() {
	Test()
	//readFileToBytes()
}
