# sgrm
simple goroutine manager

```
func myFunc2(num int) {
	fmt.Printf("num:%v\n", num)
	time.Sleep(time.Second * 1)
}
func main() {
	for i := 1; i <= 10; i++ {
		AddTask(i)
	}

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

/*
Parameters copied
*/
func AddTask(num int) {
	GRM.Add(fmt.Sprintf("task%v", num), func() { myFunc2(num) })
}
```
#full example
```
func myFunc1() {
	fmt.Println("this is test")
	time.Sleep(time.Second * 1)
}
func main() {
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

```
