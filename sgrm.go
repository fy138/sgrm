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
	IsPause     bool
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
	}

	if _, loaded := m.routines.LoadOrStore(name, &managedRoutine{Name: name}); loaded {
		return fmt.Errorf("routine with name '%s' already exists", name)
	}

	m.routines.Store(name, routine)
	m.Wg.Add(1)

	go func() {
		fmt.Printf("Task %s started\n", name)
		defer m.Wg.Done()
		for {
			select {
			case <-m.ctx.Done():
				fmt.Printf("Task '%s' stopped. Started at: %v, Last Run at: %v\n", name,
					routine.StartTime.Format("2006-01-02 15:04:05"), routine.LastRunTime.Format("2006-01-02 15:04:05"))
				return
			case <-routine.StopChn:
				fmt.Printf("Task '%s' manually stopped. Started at: %v, Last Run at: %v\n",
					name, routine.StartTime.Format("2006-01-02 15:04:05"), routine.LastRunTime.Format("2006-01-02 15:04:05"))
				return
			case <-routine.PauseChn:
				fmt.Printf("Task '%s' pasused", name)
				routine.IsPause = true
				<-routine.PauseChn
				routine.IsPause = false
				fmt.Printf("Task '%s' resume", name)

			default:
				start := time.Now().UnixMilli()
				fn()
				routine.Speed = float64((time.Now().UnixMilli() - start)) / 1000.00
				// 更新上次运行时间
				routine.LastRunTime = time.Now()
				routine.Count += 1
			}
		}
	}()
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

// StopAll 停止所有正在运行的goroutine
func (m *Manager) StopAll() {
	m.cancel()       // 通知所有goroutine其依赖的context已被取消
	close(m.stopAll) // 发送信号，辅助清理已停止的goroutine记录

	m.routines.Range(func(key, value interface{}) bool {
		m.routines.Delete(key)
		return true
	})
	m.Wg.Done()
}

func (m *Manager) PauseAll() {
	m.routines.Range(func(key, value interface{}) bool {
		value.(managedRoutine).PauseChn <- struct{}{}
		return true
	})

}

func (m *Manager) ResumeAll() {
	m.routines.Range(func(key, value interface{}) bool {
		value.(managedRoutine).PauseChn <- struct{}{}
		return true
	})
}

func (m *Manager) Stop(name string) error {
	routine, ok := m.routines.Load(name)
	if !ok {
		return fmt.Errorf("routine not exists.")
	}
	routine.(managedRoutine).StopChn <- struct{}{}
	return nil
}

func (m *Manager) Pause(name string) error {
	routine, ok := m.routines.Load(name)
	if !ok {
		return fmt.Errorf("routine not exists.")
	}
	routine.(managedRoutine).PauseChn <- struct{}{}
	return nil

}
func (m *Manager) Resume(name string) error {
	routine, ok := m.routines.Load(name)
	if !ok {
		return fmt.Errorf("routine not exists.")
	}
	routine.(managedRoutine).PauseChn <- struct{}{}
	return nil
}
