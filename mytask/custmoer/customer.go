package customer

type Customer struct {
	topic   string
	closeCh chan struct{}
}

func (c *Customer) Run() {
	for {
		task := queue.Pop(c.topic)
		if task != nil {
			c.do(task)
			continue
		}

		time.Sleep(time.Second * 1)
	}
}
