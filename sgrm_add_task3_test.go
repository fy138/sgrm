package sgrm

import (
	"fmt"
	"testing"
	"time"
)

func TestAddTask3(t *testing.T) {
	for i := 1; i <= 10; i++ {
		AddTask(i)
	}

	list := GRM.GetAllTask()
	for _, goroutine := range list {
		fmt.Printf("name: %s ->start time: %s -> last run: %s -> count %d\n",
			goroutine.Name,
			goroutine.StartTime.Format("2006-01-02 13:04:05"),
			goroutine.LastRunTime.Format("2006-01-02 13:04:05"),
			goroutine.Count,
		)
	}
	GRM.StartAll()
	time.Sleep(time.Second * 5)
	GRM.StopAll()
}

/*
Parameters copied
*/
func AddTask(num int) {
	GRM.Add(fmt.Sprintf("task%v", num), func() { myFunc2(num) })
}
