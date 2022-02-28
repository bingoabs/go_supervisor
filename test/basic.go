package test

import (
	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

type HttpWorker struct {
	Name  string
	Value map[string]string
}

func (work *HttpWorker) Get(status interface{}) (interface{}, error) {
	return status, nil
}

func (work *HttpWorker) Refresh(status interface{}) (Supervisor.IWorker, error) {
	// 在RefreshInterval为0时，不进行刷新
	return work, nil
}

func CreateHttpWorker(entry *Supervisor.Entry) Supervisor.IWorker {
	work := HttpWorker{
		Name:  entry.Name,
		Value: make(map[string]string, 5),
	}
	return &work
}

func getSupervisorOption() Supervisor.Option {
	option := Supervisor.Option{
		Name:      "test_supervisor",
		Generator: CreateHttpWorker,
		RestartRule: Supervisor.Strategy{
			TimeInterval: 10,
			Number:       3,
		},
		RefreshInterval: 5,
	}
	return option
}
