package test

import (
	"log"

	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

type HttpWorker struct {
	Name  string
	Value map[string]string
}

func (work *HttpWorker) Get(status interface{}) (interface{}, error) {
	return status, nil
}

// best to return the new object, not the origin one
func (work *HttpWorker) Refresh(status interface{}) (Supervisor.IWorker, error) {
	// 在RefreshInterval为0时，不进行刷新
	log.Println("Customer function Refresh start")
	new_worker := &HttpWorker{
		Name:  "After_refresh",
		Value: work.Value,
	}
	return new_worker, nil
}

func CreateHttpWorker(entry *Supervisor.Entry) Supervisor.IWorker {
	work := HttpWorker{
		Name:  entry.Name,
		Value: make(map[string]string, 5),
	}
	return &work
}
