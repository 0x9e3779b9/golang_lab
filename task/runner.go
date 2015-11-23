package main

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"golang_lab/task/models"
)

const (
	RUNNING = iota
	PAUSE
	STOP
	FINISH
)

type TaskUnit struct {
	source     *models.Task
	tasks      []*models.TaskCase
	records    map[int64]*models.Record
	wg         *sync.WaitGroup
	index      int32
	pauseCh    chan struct{}
	continueCh chan struct{}
	stopCh     chan struct{}
	state      int32
}

func NewTaskUnit(task *models.Task) *TaskUnit {
	return &TaskUnit{
		source:     task,
		pauseCh:    make(chan struct{}),
		stopCh:     make(chan struct{}),
		continueCh: make(chan struct{}),
		wg:         new(sync.WaitGroup),
	}
}

func (t *TaskUnit) Init() {

}

func (t *TaskUnit) pause() {
	t.pauseCh <- struct{}{}
	atomic.StoreInt32(&t.state, PAUSE)
}

func (t *TaskUnit) proceed() {
	t.continueCh <- struct{}{}
	atomic.StoreInt32(&t.state, RUNNING)
}

func (t *TaskUnit) stop() {
	t.stopCh <- struct{}{}
	atomic.StoreInt32(&t.state, STOP)
}

func (t *TaskUnit) push() {
	index := atomic.LoadInt32(&t.index)
	defer atomic.StoreInt32(&t.index, index+1)

	tc := t.tasks[index]
	rd := models.NewRecord(tc)
	t.records[rd.ID] = rd
	redis.PushTask(rd)
	ret, _ := redis.PopTaskResponse()
	rd.Result(ret)
	models.InsertRecord(rd)
}

func (t *TaskUnit) run() {
	var tag int32
	t.wg.Add(1)
	defer t.wg.Done()
LOOP:
	for {
		select {
		case <-t.pauseCh:
			atomic.StoreInt32(&tag, 1)
		case <-t.continueCh:
			atomic.StoreInt32(&tag, 0)
		case <-t.stopCh:
			break LOOP
		default:
			if atomic.LoadInt32(&tag) == 0 {
				t.push()
				if int(t.index) >= len(t.tasks) {
					atomic.StoreInt32(&t.state, FINISH)
					break LOOP
				}
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}
}

type TaskManager struct {
	pool       map[int64]*TaskUnit
	inCh       chan *models.Task
	stop       chan struct{}
	pauseCh    chan int64
	continueCh chan int64
	stopCh     chan int64
}

var defaultTaskMgr *TaskManager

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
		t  *models.Tasl
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
