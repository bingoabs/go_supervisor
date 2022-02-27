package test

import (
	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

type HttpWorker struct {
	Name  string
	Value string
}


func (work *HttpWorker) Init(entry *Supervisor.Entry) Supervisor.IWorker {


func (work *HttpWorker) Get(status interface{}) (interface{}, error) {

}

func (work *HttpWorker) Put(status interface{}) (Supervisor.IWorker, error) {

}

func (work *HttpWorker) Refresh(status interface{}) (Supervisor.IWorker, error) {

}

func CreateHttpWorker(entry *Supervisor.Entry) Supervisor.IWorker {
	
}

func getSupervisorOption() Supervisor.Option {

	// const option Supervisor.Option = Supervisor.Option{
	// 	Name:   "test_supervisor",
	// 	Worker: HttpWorker{},
	// 	RestartRule: Supervisor.Strategy{
	// 		TimeInterval: 10,
	// 		Number:       3,
	// 	},
	// 	RefreshInterval: 0,
	// }
}
https://zhuanlan.zhihu.com/p/95473929
https://www.zhihu.com/question/30461290/answer/2366427697

