package sgrm

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// simple goroutine manager

var GRM *Manager

func init() {
	GRM = NewManager()
}

type managedRoutine struct {
	Name        string
	StopChn     chan struct{}
	PauseChn    chan struct{}
	StartTime   time.Time
	LastRunTime time.Time
	Count       int
	Speed       float64 //seconds
	IsPaused    bool
	Istarted    bool //是否已经启动
	Task        func()
}

type Manager struct {
	//mu       sync.Mutex
	routines *sync.Map
	stopAll  chan struct{}
	ctx      context.Context
	cancel   context.CancelFunc
	Wg       *sync.WaitGroup
}

func NewManager() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		routines: &sync.Map{},
		stopAll:  make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
		Wg:       &sync.WaitGroup{},
	}
}

// Add 添加一个新的goroutine到管理器中
func (m *Manager) Add(name string, fn func()) error {
	routine := &managedRoutine{
		Name:        name,
		StopChn:     make(chan struct{}),
		PauseChn:    make(chan struct{}),
		LastRunTime: time.Now(),
		StartTime:   time.Now(),
		Task:        fn,
	}

	if _, loaded := m.routines.Load(name); loaded {
		return fmt.Errorf("routine with name '%s' already exists", name)
	}

	m.routines.Store(name, routine)

	return nil
}

func (m *Manager) GetAllTask() []*managedRoutine {
	list := []*managedRoutine{}
	m.routines.Range(func(key, value any) bool {
		r := value.(*managedRoutine)
		list = append(list, r)
		return true
	})
	return list
}

/*
*
run all  task
*/
func (m *Manager) StartAll() {
	m.routines.Range(func(key, value interface{}) bool {
		//m.routines.Delete(key)
		routine := value.(*managedRoutine)
		if routine.Istarted {
			return true
		}
		m._start(routine)
		return true
	})
	//m.Wg.Wait()
}

/*
*
run a  task
*/
func (m *Manager) Start(name string) error {
	mr, err := m.GetTaskByName(name)
	if err != nil {
		return err
	}
	if mr.Istarted {
		return fmt.Errorf("taks %v was started.", name)
	}
	m._start(mr)
	return nil
}

/*
*
run a  task
*/
func (m *Manager) _start(routine *managedRoutine) {
	m.Wg.Add(1)
	routine.Istarted = true
	go func() {
		fmt.Printf("Task %s started\n", routine.Name)
		defer m.Wg.Done()
		for {
			select {
			case <-m.ctx.Done():
				fmt.Printf("Task '%s' stopped. Started at: %v, Last Run at: %v\n", routine.Name,
					routine.StartTime.Format("2006-01-02 15:04:05"), routine.LastRunTime.Format("2006-01-02 15:04:05"))
				return
			case <-routine.StopChn:
				fmt.Printf("Task '%s' manually stopped. Started at: %v, Last Run at: %v\n",
					routine.Name, routine.StartTime.Format("2006-01-02 15:04:05"), routine.LastRunTime.Format("2006-01-02 15:04:05"))
				return
			case <-routine.PauseChn:
				fmt.Printf("Task '%s' pasused\n", routine.Name)
				routine.IsPaused = true
				<-routine.PauseChn
				routine.IsPaused = false
				fmt.Printf("Task '%s' resumed\n", routine.Name)

			default:
				start := time.Now().UnixMilli()
				routine.Task()
				routine.Speed = float64((time.Now().UnixMilli() - start)) / 1000.00
				// 更新上次运行时间
				routine.LastRunTime = time.Now()
				routine.Count += 1
			}
		}
	}()
}

/*
*
get a task by name
*/
func (m *Manager) GetTaskByName(name string) (*managedRoutine, error) {
	routine, ok := m.routines.Load(name)
	if !ok {
		return nil, fmt.Errorf("routine with name '%s' not exists", name)
	}
	return routine.(*managedRoutine), nil
}

/*
*
stop  all task
*/
func (m *Manager) StopAll() {
	m.cancel()
	close(m.stopAll)

	m.routines.Range(func(key, value interface{}) bool {
		m.routines.Delete(key)
		return true
	})
	m.Wg.Wait()
}

/*
*
stop all task
*/
func (m *Manager) Stop(name string) error {
	mr, err := m.GetTaskByName(name)
	if err != nil {
		return err
	}
	mr.StopChn <- struct{}{}
	m.routines.Delete(name)
	return nil
}

/*
*
pause all task
*/
func (m *Manager) PauseAll() {
	m.routines.Range(func(key, value interface{}) bool {
		value.(*managedRoutine).PauseChn <- struct{}{}
		return true
	})

}

/*
*
pause a task
*/
func (m *Manager) Pause(name string) error {
	mr, err := m.GetTaskByName(name)
	if err != nil {
		return err
	}
	mr.PauseChn <- struct{}{}
	return nil
}

/*
*
resume all task
*/
func (m *Manager) ResumeAll() {
	m.routines.Range(func(key, value interface{}) bool {
		value.(*managedRoutine).PauseChn <- struct{}{}
		return true
	})
}

/*
*
resume a task
*/
func (m *Manager) Resume(name string) error {
	mr, err := m.GetTaskByName(name)
	if err != nil {
		return err
	}
	mr.PauseChn <- struct{}{}
	return nil
}
