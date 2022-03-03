package test

import (
	"log"
	"testing"
	"time"

	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

const ENTRY_NAME string = "entry_1"

func TestCreate(t *testing.T) {
	log.Println("TestGetPush start")
	option := getSupervisorOption()
	supervisor := Supervisor.CreateSupervisor(option)
	_ = Supervisor.GetWorker(supervisor, "entry_1")
	_ = Supervisor.GetWorker(supervisor, "entry_2")
}

func TestWorkerGet(t *testing.T) {
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

func TestWorkerClose(t *testing.T) {
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
	option := getSupervisorOption()
	option.RefreshInterval = 3
	supervisor := Supervisor.CreateSupervisor(option)
	entry := Supervisor.GetWorker(supervisor, ENTRY_NAME)
	time.Sleep(6 * time.Second)
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

func TestWorkerRefreshCrash(t *testing.T) {
	option := getPanicWorkerSupervisorOption()
	supervisor := Supervisor.CreateSupervisor(option)
	entry := Supervisor.GetWorker(supervisor, ENTRY_NAME)
	time.Sleep(10 * time.Second)
	response := Supervisor.GetFromEntry(entry, ENTRY_NAME)
	log.Println("response, ", response)
	if response.Err == nil {
		t.Errorf("Worker Panic goroutine can't Get error")
	}
	err, ok := response.Err.(Supervisor.EntryPanicError)
	if !ok {
		t.Errorf("Worker Panic goroutine can't Get Panic Error")
	}
	log.Printf("Worker Panic goroutine Get Panic error: %v\n", err)
}
