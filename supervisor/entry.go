package supervisor

import (
	"log"
	"sync"
	"time"
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

func Get(entry *Entry, message interface{}) interface{} {
	return nil
}

func Put(entry *Entry, message interface{}) {
	return
}

func Close(entry *Entry) {
	log.Println("Entry Close start ", entry)
	if entry.e_closed {
		return
	}
	entry.e_lock.Lock()
	if !entry.e_closed {
		close(entry.Mq)
		entry.e_closed = true
	}
	entry.e_lock.Unlock()
	log.Println("Entry Close end: ", entry)
}

func send_refersh_signal(entry *Entry, interval int) {
	log.Println("Entry send_refersh_signal start with interval: ", interval)
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	select {
	case <-ticker.C:
		message := WorkerReceiveMessage{
			MessageType: WORKER_MESSAGE_REFRESH,
		}
		send_message(entry, message)
	}
}

func send_message(entry *Entry, message WorkerReceiveMessage) error {
	发送信息
	TODO
	考虑对所有发出的信息，worker都进行回复
}
