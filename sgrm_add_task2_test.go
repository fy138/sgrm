package sgrm

import (
	"fmt"
	"testing"
	"time"
)

func myFunc2(num int) {
	fmt.Printf("num:%v\n", num)
	time.Sleep(time.Second * 1)
}

func TestAddTask2(t *testing.T) {
	for i := 1; i <= 10; i++ {
		x := i //must be new value
		GRM.Add(fmt.Sprintf("task%v", i), func() {
			//fmt.Println(i)
			myFunc2(x)
		})
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
