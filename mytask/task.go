package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang_lab/mytask/message"
	"golang_lab/mytask/notify"
	"golang_lab/mytask/queue"
)

const (
	RUNNING = iota
	PAUSE
	STOP
	FINISH
)

type History struct {
	ID       int64 `json:"id"`
	TaskId   int64 `json:"task_id"`
	Cost     int
	Start    int64
	End      int64
	TaskConf string
	Params   map[string]string
	Owner    string
	State    int64
	LogPath  string
	Result   int64
}

func (h *History) encode() string {
	byt, _ := json.Marshal(h)
	return string(byt)
}

type Task struct {
	ID         int64
	hs         map[int64]*History
	wg         *sync.WaitGroup
	pauseCh    chan struct{}
	continueCh chan struct{}
	stopCh     chan struct{}
	state      int64
	cnt        int64
	index      int64
}

func NewTask(taskID, hsNum int64) *Task {
	var i int64
	hs := make(map[int64]*History)
	for i = 1; i <= hsNum; i++ {
		hs[i] = &History{ID: i, TaskId: taskID}
	}
	return &Task{
		ID:         taskID,
		hs:         hs,
		cnt:        hsNum,
		pauseCh:    make(chan struct{}, 2),
		stopCh:     make(chan struct{}, 2),
		continueCh: make(chan struct{}, 2),
		wg:         new(sync.WaitGroup),
	}
}

func (t *Task) Pause() {
	t.pauseCh <- struct{}{}
	atomic.StoreInt64(&t.state, PAUSE)
}

func (t *Task) Continue() {
	t.continueCh <- struct{}{}
	atomic.StoreInt64(&t.state, RUNNING)
}

func (t *Task) Stop() {
	t.stopCh <- struct{}{}
	atomic.StoreInt64(&t.state, STOP)
}

func (t *Task) do(id int64) {
	// push
	topic := fmt.Sprintf("task_%d", t.ID)
	err := queue.Push(topic, t.hs[t.index+1].encode())
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// get logPath
	topic = fmt.Sprintf("task_%d_resp", t.ID)
	for {
		v, _ := queue.Pop(topic)
		if v != nil {
			fmt.Printf("Got ack response %v from %s\n", v, topic)
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}

	// get finish
	for {
		v, _ := queue.Pop(topic)
		if v != nil {
			fmt.Printf("Got finish response %v from %s\n", v, topic)
			break
		} else {
			time.Sleep(time.Second * 1)
		}
	}
	fmt.Printf("task_%d -- %d finished!\n", t.ID, id)
}

func (t *Task) Run() {
	var tag int64
	t.wg.Add(1)
	defer t.wg.Done()
LOOP:
	for {
		select {
		case <-t.pauseCh:
			atomic.StoreInt64(&tag, 1)
		case <-t.continueCh:
			atomic.StoreInt64(&tag, 0)
		case <-t.stopCh:
			break LOOP
		default:
			if atomic.LoadInt64(&tag) == 0 {
				t.do(t.index + 1)
				t.index++
				if t.index >= t.cnt {
					atomic.StoreInt64(&t.state, FINISH)
					break LOOP
				}
			} else {
				time.Sleep(time.Second * 1)
			}
		}
	}
}

type TaskManager struct {
	pool      map[int64]*Task
	notifier  *notify.FSnotifyManager
	msgCenter *message.MessageCenter

	lock *sync.Mutex
}

func NewTaskManager() *TaskManager {
	mc := message.NewMessageCenter()
	return &TaskManager{
		pool:      make(map[int64]*Task),
		notifier:  notify.NewFSnotifyManager(8, mc.Input),
		msgCenter: mc,
		lock:      new(sync.Mutex),
	}
}

func (tm *TaskManager) Prepared() {
	go tm.notifier.Run()
	go tm.msgCenter.Run()
}

func (tm *TaskManager) Start(id int64) (err error) {
	fmt.Printf("Start start %d\n", id)
	defer fmt.Printf("Start finished %d\n", id)
	tm.lock.Lock()
	if _, ok := tm.pool[id]; ok {
		if tm.pool[id].state != FINISH {
			err = fmt.Errorf("already running")
			return
		}
		delete(tm.pool, id)
	}
	tm.pool[id] = NewTask(id, 10)
	tm.notifier.Push(fmt.Sprintf("/tmp/task_%d", id))
	go tm.pool[id].Run()
	tm.lock.Unlock()
	return
}

func (tm *TaskManager) Stop(id int64) (err error) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if _, ok := tm.pool[id]; ok {
		err = fmt.Errorf("not start")
		return
	}
	if tm.pool[id].state != RUNNING {
		err = fmt.Errorf("not running")
		return
	}
	tm.pool[id].Stop()
	return
}

func (tm *TaskManager) Pause(id int64) (err error) {
	fmt.Printf("pause start %d\n", id)
	defer fmt.Printf("pause finished %d\n", id)
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if _, ok := tm.pool[id]; !ok {
		err = fmt.Errorf("not start")
		return
	}
	if tm.pool[id].state != RUNNING {
		err = fmt.Errorf("not running")
		return
	}
	tm.pool[id].Pause()
	return
}

func (tm *TaskManager) Continue(id int64) (err error) {
	fmt.Printf("Continue start %d\n", id)
	defer fmt.Printf("Continue finished %d\n", id)
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if _, ok := tm.pool[id]; !ok {
		err = fmt.Errorf("not start")
		return
	}
	if tm.pool[id].state != PAUSE {
		err = fmt.Errorf("not pause")
		return
	}
	tm.pool[id].Continue()
	return
}

func main() {
	err := queue.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tm := NewTaskManager()
	tm.Prepared()

	tm.Start(1)
	tm.Start(2)

	time.Sleep(time.Second * 2)

	tm.Pause(1)
	tm.Pause(2)

	time.Sleep(time.Second * 2)

	tm.Continue(1)
	tm.Continue(2)

	tk := time.NewTicker(time.Second * 4)
	defer tk.Stop()
	for {
		select {
		case <-tk.C:
			//fmt.Println()
		}
	}
}
