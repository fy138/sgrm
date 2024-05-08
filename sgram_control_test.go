package sgrm

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestAddTask4(t *testing.T) {
	GRM.Add("task1", myFunc1)
	GRM.Add("task2", myFunc1)

	time.Sleep(time.Second * 5)

	GRM.PauseAll() //pause all

	list := GRM.GetAllTask()
	for _, goroutine := range list {
		fmt.Printf("%+v\n", goroutine)
	}
	time.Sleep(time.Second * 5)
	GRM.ResumeAll() //resume all

	time.Sleep(time.Second * 5)
	list = GRM.GetAllTask()
	for _, goroutine := range list {
		fmt.Printf("%+v\n", goroutine)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	GRM.StopAll()
}
