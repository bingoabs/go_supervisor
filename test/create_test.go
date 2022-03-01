package test

import (
	"log"
	"testing"

	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

const ENTRY_NAME string = "entry_1"

func aTestCreate(t *testing.T) {
	log.Println("TestGetPush start")
	option := getSupervisorOption()
	supervisor := Supervisor.CreateSupervisor(option)
	_ = Supervisor.GetWorker(supervisor, "entry_1")
	_ = Supervisor.GetWorker(supervisor, "entry_2")
}

func aTestWorkerGet(t *testing.T) {
	option := getSupervisorOption()
	supervisor := Supervisor.CreateSupervisor(option)
	entry := Supervisor.GetWorker(supervisor, ENTRY_NAME)
	response := Supervisor.GetFromEntry(entry, ENTRY_NAME)
	if response.Err != nil {
		t.Errorf("Worker Get receive error: %v", response.Err)
	}
	data, ok := response.Data.(string)
	if !ok {
		t.Errorf("Worker Get response data: %v", response.Data)
	}
	if data != ENTRY_NAME {
		t.Errorf("Worker Get data not match: %v", data)
	}
}

func aTestWorkerClose(t *testing.T) {
	option := getSupervisorOption()
	supervisor := Supervisor.CreateSupervisor(option)
	entry := Supervisor.GetWorker(supervisor, ENTRY_NAME)
	Supervisor.CloseEntry(entry)
	response := Supervisor.GetFromEntry(entry, ENTRY_NAME)
	if response.Err == nil {
		t.Errorf("Closed Worker Get receive")
	}
	err, ok := response.Err.(Supervisor.EntryPanicError)
	if !ok {
		t.Errorf("Closed Worker Get unknow error: %v", err)
	}
}

func TestWorkerRefresh(t *testing.T) {
	// TODO
	option := getSupervisorOption()
	supervisor := Supervisor.CreateSupervisor(option)
	entry := Supervisor.GetWorker(supervisor, ENTRY_NAME)
	response := Supervisor.GetFromEntry(entry, ENTRY_NAME)
	if response.Err != nil {
		t.Errorf("Worker Get receive error: %v", response.Err)
	}
	data, ok := response.Data.(string)
	if !ok {
		t.Errorf("Worker Get response data: %v", response.Data)
	}
	if data != ENTRY_NAME {
		t.Errorf("Worker Get data not match: %v", data)
	}
}
