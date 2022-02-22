package supervisor

import (
	"sync"
)

type Entry struct {
	Name     string
	Mq       chan WorkerReceiveMessage
	e_closed bool
	e_lock   sync.Mutex
}

type EntryStatic struct {
	Entry            *Entry
	Panic_timestamps []int64
}

func (entry *Entry) Sync(message interface{}) interface{} {
	return nil
}

func (entry *Entry) Async(message interface{}) {
	return
}
