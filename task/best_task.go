package main

type TaskManager struct {
	pool       map[int64]*TaskUnit
	inCh       chan *models.Task
	stop       chan struct{}
	pauseCh    chan int64
	continueCh chan int64
	stopCh     chan int64
}

var defaultTaskMgr *TaskManager

func Init() {
	defaultTaskMgr = &TaskManager{
		pool:       make(map[int64]*TaskUnit),
		inCh:       make(chan *models.Task, 24),
		stop:       make(chan struct{}),
		pauseCh:    make(chan int64, 24),
		continueCh: make(chan int64, 24),
		stopCh:     make(chan int64, 24),
	}
}

func AcceptTask(t *models.Task) (err error) {
	if _, ok := defaultTaskMgr.pool[t.ID]; ok {
		err = errors.New("Existed!")
		return
	}

	defaultTaskMgr.inCh <- t
	return
}

func PauseTask(id int64) (err error) {
	if _, ok := defaultTaskMgr.pool[id]; !ok {
		err = errors.New("NotExisted!")
		return
	}
	if atomic.LoadInt32(&defaultTaskMgr.pool[id].state) != RUNNING {
		err = errors.New("NotRunning!")
		return
	}

	defaultTaskMgr.pauseCh <- id
	return
}

func StopTask(id int64) (err error) {
	if _, ok := defaultTaskMgr.pool[id]; !ok {
		err = errors.New("NotExisted!")
		return
	}
	if atomic.LoadInt32(&defaultTaskMgr.pool[id].state) >= STOP {
		err = errors.New("Has Stopped or finished!")
		return
	}
	defaultTaskMgr.stopCh <- id
	return
}

func ProceedTask(id int64) (err error) {
	if _, ok := defaultTaskMgr.pool[id]; !ok {
		err = errors.New("NotExisted!")
		return
	}
	if atomic.LoadInt32(&defaultTaskMgr.pool[id].state) != PAUSE {
		err = errors.New("Not Pause!")
		return
	}
	defaultTaskMgr.continueCh <- id
	return
}

func (m *TaskManager) ReadLoop() {
	var (
		t  *models.Task
		id int64
	)

	for {
		select {
		case t = <-m.inCh:
			if _, ok := defaultTaskMgr.pool[t.ID]; !ok {
				defaultTaskMgr.pool[t.ID] = NewTaskUnit(t)
				go defaultTaskMgr.pool[t.ID].run()
			}

		case id = <-m.pauseCh:
			if _, ok := defaultTaskMgr.pool[t.ID]; ok {
				defaultTaskMgr.pool[t.ID].pause()
			}
		case id = <-m.stopCh:
			if _, ok := defaultTaskMgr.pool[t.ID]; ok {
				defaultTaskMgr.pool[t.ID].stop()
			}
		case id = <-m.continueCh:
			if _, ok := defaultTaskMgr.pool[t.ID]; ok {
				defaultTaskMgr.pool[t.ID].proceed()
			}
		case <-m.stop:
			return
		}
	}
}
