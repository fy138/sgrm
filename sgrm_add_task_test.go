package sgrm

import (
	"fmt"
	"testing"
	"time"
)

func myFunc1() {
	fmt.Println("this is test1")
	time.Sleep(time.Second * 1)
}
func myFunc1a() {
	fmt.Println("this is test2")
	time.Sleep(time.Second * 1)
}
func TestAddTask1(t *testing.T) {
	fmt.Println("添加任务")
	GRM.Add("task1", myFunc1)
	GRM.Add("task2", myFunc1a)

	//time.Sleep(time.Second * 5)

	list := GRM.GetAllTask()
	for _, goroutine := range list {
		fmt.Printf("%+v\n", goroutine)
	}
	//time.Sleep(time.Second * 5)
	fmt.Println("开始任务")
	GRM.StartAll()
	time.Sleep(time.Second * 5)
	fmt.Println("停止任务")
	GRM.StopAll()
}
func TestAddTask1a(t *testing.T) {
	fmt.Println("添加任务")
	GRM.Add("task1", myFunc1a)
	if err := GRM.Add("task1", myFunc1); err != nil {
		fmt.Println(err)
	}
}

func TestAddTask1b(t *testing.T) {
	fmt.Println("添加任务")
	GRM.Add("task1", myFunc1)
	GRM.Add("task2", myFunc1a)

	fmt.Println("开始任务")
	GRM.Start("task1")

	time.Sleep(time.Second * 5)
	list := GRM.GetAllTask()
	for _, goroutine := range list {
		fmt.Printf("%+v\n", goroutine)
	}
	//time.Sleep(time.Second * 5)
	fmt.Println("停止任务")
	GRM.StopAll()
}
