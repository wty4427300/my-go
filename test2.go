package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

var (
//Max_Num = os.Getenv("MAX_NUM")
MaxWorker = runtime.NumCPU()//cpu数
MaxQueue = 1000//队列容量1000
)

type Serload struct {
	pri string
}
//正式使用应该将job改成接口
type Job struct {
	serload Serload
}

var JobQueue chan Job

type Worker struct {
	WorkerPool chan chan Job//工作池
	JobChannel chan Job//工作
	Quit chan bool//退出
}

func NewWorker(workPool chan chan Job) Worker {
	return Worker {
		WorkerPool:workPool,
		JobChannel:make(chan Job),
		Quit:make(chan bool),}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
				case job:= <- w.JobChannel:
				// excute job
				fmt.Println(job.serload.pri)
				case <- w.Quit:
				return}
			}
		}()
}

func (w Worker) Stop() {
	go func() {
		w.Quit <- true
		}()
}

type Dispatcher struct {//调度器
	MaxWorkers int
	WorkerPool chan chan Job
	Quit chan bool
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{MaxWorkers:maxWorkers, WorkerPool:pool, Quit:make(chan bool)}
}

func (d *Dispatcher) Run() {
	for i:=0; i<d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
		}

	go d.Dispatch()
}

func (d *Dispatcher) Stop() {
	go func() {
		d.Quit <- true
		}()
}

func (d *Dispatcher) Dispatch() {
	for {
		select {
			case job:=<- JobQueue:
			go func(job Job) {
				jobChannel := <- d.WorkerPool
				jobChannel <- job
				}(job)
			case <- d.Quit:
			return
			}
		}
}

func entry(res http.ResponseWriter, req *http.Request) {
	// fetch job
	work := Job{serload:Serload{pri:"Just do it"}}
	JobQueue <- work
	fmt.Fprintf(res, "滚滚滚滚")
}

func init() {
	runtime.GOMAXPROCS(MaxWorker)
	JobQueue = make(chan Job, MaxQueue)
	dispatcher := NewDispatcher(MaxWorker)
	dispatcher.Run()
}

func main() {
	Port := "8086"
	IsHttp := true
	arg_num := len(os.Args)
	fmt.Println(arg_num)
	if 2<=arg_num {
		Port = os.Args[1]
		}
	if 3<=arg_num {
		if os.Args[2]=="true" {
			IsHttp = true
			} else {
			IsHttp = false
			}
		}
	fmt.Printf("server is http %t\n", IsHttp)
	fmt.Println("server listens at ", Port)

	http.HandleFunc("/", entry)

	var err error
	if IsHttp {
		err = http.ListenAndServe(":"+Port, nil)
		} else {
		err = http.ListenAndServeTLS(":"+Port, "server.crt", "server.key", nil)
		}
	if err != nil {
		fmt.Println("Server failure /// ", err)
		}

	fmt.Println("quit")
}