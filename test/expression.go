package test

import (
	"fmt"
)

type StackNode struct {
	Data interface{}
	next *StackNode
}

type LinkStack struct {
	top   *StackNode
	Count int
}

func (this *LinkStack) Init() {
	this.top = nil
	this.Count = 0
}

func (this *LinkStack) Push(data interface{}) {
	var node *StackNode = new(StackNode)
	node.Data = data
	node.next = this.top
	this.top = node
	this.Count++
}

func (this *LinkStack) Pop() interface{} {
	if this.top == nil {
		return nil
	}
	returnData := this.top.Data
	this.top = this.top.next
	this.Count--
	return returnData
}

func (this *LinkStack) LookTop() interface{} {
	if this.top == nil {
		return nil
	}
	return this.top.Data
}

func changeMidExp2PostExp(src string) {
	dst := make([]byte, 0)
	st := new(LinkStack)
	st.Init()
	for i, _ := range src {
		switch src[i] {
		case byte('('):
			st.Push(byte('('))
		case byte(')'):
			for {
				e := st.Pop()
				if e == nil {
					break
				}
				if e.(byte) == '(' {
					break
				}
				dst = append(dst, e.(byte))
			}
		case byte('|'), byte('&'):
			if st.LookTop().(byte) == '-' {
				st.Push(byte(src[i]))
			} else {
				for {
					e := st.Pop()
					if e == nil {
						break
					}
					if e.(byte) == '-' || e.(byte) == '(' {
						break
					}
					dst = append(dst, e.(byte))
				}
				st.Push(byte(src[i]))
			}
		case byte('+'), byte('-'):
			for {
				e := st.Pop()
				if e == nil {
					break
				}
				if e.(byte) == '(' || e.(byte) == ')' {
					break
				}
				dst = append(dst, e.(byte))
			}
			st.Push(byte(src[i]))
		default:
			dst = append(dst, byte(src[i]))
		}

	}
	fmt.Println(string(dst))
}

/*
func main() {
	src := "(a-d)&(b-c)"
	changeMidExp2PostExp(src)

}*/
