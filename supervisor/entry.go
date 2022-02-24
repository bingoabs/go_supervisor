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

type EntryStatic struct {
	Entry            *Entry
	Panic_timestamps []int64
}

func Get(entry *Entry, message interface{}) WorkerResponseMessage {
	msg := WorkerReceiveMessage{
		MessageType: WORKER_MESSAGE_GET,
		Data:        message,
		Mq:          make(chan WorkerResponseMessage, 1),
	}
	return send_message(entry, msg)
}

// 基于作为缓存的功能，put的作用不明显，尤其不适应节点间重排后的状态不稳定状态
// func Put(entry *Entry, message interface{}) WorkerResponseMessage {
// 	msg := WorkerReceiveMessage{
// 		MessageType: WORKER_MESSAGE_PUT,
// 		Data:        message,
// 		Mq:          make(chan WorkerResponseMessage, 1),
// 	}
// 	return send_message(entry, msg)
// }

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

func send_message(entry *Entry, message WorkerReceiveMessage) WorkerResponseMessage {
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
