package supervisor

import (
	"log"
	"sync"
)

type Entry struct {
	Name       string
	Mq         chan WorkerReceiveMessage
	Created_At int64
	e_closed   bool
	e_lock     sync.Mutex
}

func Get(entry *Entry, message interface{}) WorkerResponseMessage {
	msg := WorkerReceiveMessage{
		MessageType: WORKER_MESSAGE_GET,
		Data:        message,
		Mq:          make(chan WorkerResponseMessage, 1),
	}
	return send_message_to_worker(entry, msg)
}

// TODO 可以使用sync/atomic进行优化
func Close(entry *Entry) {
	log.Println("Entry Close start ", entry)
	if entry.e_closed {
		return
	}
	entry.e_lock.Lock()
	if !entry.e_closed {
		entry.e_closed = true
		close(entry.Mq)
	}
	entry.e_lock.Unlock()
	log.Println("Entry Close end: ", entry)
}

func send_message_to_worker(entry *Entry, message WorkerReceiveMessage) WorkerResponseMessage {
	if entry.e_closed {
		return WorkerResponseMessage{
			Err:  EntryPanicError{},
			Data: nil,
		}
	}
	entry.e_lock.Lock()
	if !entry.e_closed {
		entry.Mq <- message
	}
	entry.e_lock.Unlock()
	result := <-message.Mq
	return result
}
