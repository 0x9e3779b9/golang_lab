package main

type Node struct {
	row     int
	col     int
	Element *Cell
	down    *Node
	right   *Node
}

type CrossLink struct {
	ColHeader []*Node
	RowHeader []*Node
	col       int
	row       int
	nz        int
}

func NewCrossLink() *CrossLink {
	return &CrossLink{
		ColHeader: []*Node{&Node{}},
		RowHeader: []*Node{&Node{}},
	}
}

func (c *CrossLink) extend(row, col int) {
	lenRow := len(c.RowHeader)
	lenCol := len(c.ColHeader)
	if lenRow-1 < row {
		for i := lenRow - 1; i <= row; i++ {
			c.RowHeader = append(c.RowHeader, &Node{})
		}
		c.row = row
	}
	if lenCol-1 < col {
		for i := lenCol - 1; i <= col; i++ {
			c.ColHeader = append(c.ColHeader, &Node{})
		}
		c.col = col
	}
}

func (c *CrossLink) Insert(row, col, v int) {
	c.extend(row, col)
	node := &Node{
		row:     row,
		col:     col,
		Element: &Cell{Val: v},
	}

	rPointer := c.RowHeader[row]
	cPointer := c.ColHeader[col]
	for {
		if rPointer.right == nil {
			rPointer.right = node
			break
		}
		if rPointer.right.col > col {
			node.right = rPointer.right
			rPointer.right = node
			break
		}
		rPointer = rPointer.right
	}

	for {
		if cPointer.down == nil {
			cPointer.down = node
			break
		}
		if cPointer.down.row > row {
			node.down = cPointer.down
			cPointer.down = node
			break
		}
		cPointer = cPointer.down
	}
}

func (c *CrossLink) Get(row, col int) {
	rPointer := c.RowHeader[row]
	for {
		if rPointer.right == nil {
			break
		}

		if rPointer.right.col == col {
			fmt.Println(rPointer.right.Element.Val)
		}

		rPointer = rPointer.right
	}
}

func (c *CrossLink) GetRow(row int) *Node {
	return c.RowHeader[row]
}

func (c *CrossLink) GetCol(col int) *Node {
	return c.ColHeader[col]
}

func main() {
	cl := NewCrossLink()
	cl.Insert(0, 1, 2012)
	cl.Insert(0, 0, 2011)
	cl.Insert(2, 1, 2013)
	cl.Insert(5, 5, 2014)
	cl.Get(0, 1)
	cl.Get(0, 0)
	cl.Get(2, 1)
	cl.Get(5, 5)
	row0 := cl.GetRow(0).right
	for {
		if row0 == nil {
			break
		}
		fmt.Printf("(%d,%d)-->%d\n", 0, row0.col, row0.data.Val)
		row0 = row0.right
	}
}
