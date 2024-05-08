package sgrm

import (
	"fmt"
	"testing"
	"time"
)

func myFunc1() {
	fmt.Println("this is test")
	time.Sleep(time.Second * 1)
}

func TestAddTask1(t *testing.T) {
	GRM.Add("task1", myFunc1)
	GRM.Add("task2", myFunc1)

	time.Sleep(time.Second * 5)

	list := GRM.GetAllTask()
	for _, goroutine := range list {
		fmt.Printf("name: %s ->start time: %s -> last run: %s -> count %d\n",
			goroutine.Name,
			goroutine.StartTime.Format("2006-01-02 13:04:05"),
			goroutine.LastRunTime.Format("2006-01-02 13:04:05"),
			goroutine.Count,
		)
	}

	GRM.StopAll()
}
