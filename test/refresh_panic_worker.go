package test

import (
	"log"

	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

type RefreshPanicHttpWorker struct {
	Name  string
	Value map[string]string
}

func (work *RefreshPanicHttpWorker) Get(status interface{}) (interface{}, error) {
	return status, nil
}

// best to return the new object, not the origin one
func (work *RefreshPanicHttpWorker) Refresh(status interface{}) (Supervisor.IWorker, error) {
	log.Panicln("Customer function Refresh panic")
	log.Println("Customer function Refresh start")
	new_worker := &RefreshPanicHttpWorker{
		Name:  "After_refresh",
		Value: work.Value,
	}
	return new_worker, nil
}

func CreateRefreshPanicHttpWorker(entry *Supervisor.Entry) Supervisor.IWorker {
	work := RefreshPanicHttpWorker{
		Name:  entry.Name,
		Value: make(map[string]string, 5),
	}
	return &work
}
