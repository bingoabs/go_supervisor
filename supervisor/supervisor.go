package supervisor

import (
	"log"
	"time"
)

type Option struct {
	Name            string
	WorkerGenerator WorkerGenerator
	RestartRule     Strategy
	RefreshInterval int
}

// 协程树根节点
type Supervisor struct {
	Name             string                        // supervisor的名称
	generator        WorkerGenrator                // worker的具体执行对象
	listen_mq        chan SupervisorReceiveMessage // supervisor监听的mq
	restart_rule     Strategy                      //worker的重启策略
	refresh_interval int                           // worker的更新时间间隔
}

/*
supervisor只管理以下事务:
1. 开启goroutine
2. 关闭goroutine
3. 监听goroutine的失效情况
4. 在goroutine故障次数超过限制时，关闭entry，调用方将发现entry关闭了，
	从而重新请求tracker，而tracker再调用supervisor对象
*/

func create_supervisor(option Option) chan SupervisorReceiveMessage {
	listen_mq := make(chan SupervisorReceiveMessage, SUPERVISOR_MQ_LENGTH)
	monitor := &Supervisor{
		Name:             option.Name,
		generator:        option.WorkerGenerator,
		listen_mq:        listen_mq,
		restart_rule:     option.RestartRule,
		refresh_interval: option.RefreshInterval,
	}
	go start_supervisor(monitor)
	return listen_mq
}

func send_message_to_supervisor(mq chan SupervisorReceiveMessage, message SupervisorReceiveMessage) *Entry {
	mq <- message
	result := <-message.Mq
	return result
}

func start_supervisor(monitor *Supervisor) {
	// 构建一个dict，用于记录奔溃的goroutine的具体统计数据，比如时间和次数
	log.Println("Supervisor start_supervisor start")
	entrys := make(map[string]*Entry)
	for {
		message, ok := <-monitor.listen_mq
		log.Println("Supervisor start_supervisor receive message")
		if !ok {
			log.Panic("Supervisor start_supervisor receive invalid message")
			return
		}
		// create, remove 由tracker调用, 而panic由定制worker调用
		if message.MessageType == SUPERVISOR_CREATE_EVENT {
			log.Println("Supervisor start_supervisor receive CREATE_EVENT")
			_, ok := entrys[message.EntryName]
			// 存在则返回, 否则新建
			if !ok {
				worker_mq := make(chan WorkerReceiveMessage, WORKER_MQ_LENGTH)
				entry := &Entry{
					Name:       message.EntryName,
					Mq:         worker_mq,
					Created_At: time.Now().Unix(),
					e_closed:   false,
				}
				entrys[message.EntryName] = entry
				go start_autoupdate_worker(monitor, entry, monitor.generator)
			}
			message.Mq <- entrys[message.EntryName]
		} else if message.MessageType == SUPERVISOR_REMOVE_EVENT {
			log.Println("Supervisor start_supervisor receive REMOVE ROUTINE action")
			_, ok := entrys[message.EntryName]
			if ok {
				go close_worker(monitor, entrys[message.EntryName])
				delete(entrys, message.EntryName)
				message.Mq <- nil
			} else {
				log.Panic("Supervisor start_supervisor receive unknown REMOVE_EVENT: ", message)
			}
		} else if message.MessageType == SUPERVISOR_DOWN_EVENT {
			// TODO 可以考虑记录更详细的信息，用于重建该entry
			log.Println("Supervisor start_supervisor receive DOWN event from: ", message)
			_, ok := entrys[message.EntryName]
			if ok {
				go close_worker(monitor, entrys[message.EntryName])
				delete(entrys, message.EntryName)
				message.Mq <- nil
			} else {
				log.Panic("Supervisor start_supervisor receive unknow Down event: ", message)
			}
		} else {
			log.Panic("Supervisor start_supervisor receive unknow EVENT: ", message)
		}
	}
}
