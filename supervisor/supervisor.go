package supervisor

import (
	"log"
)

type Option struct {
	Name        string
	Tracker_mq  chan WorkerResponseMessage
	Template    IWorker
	RestartRule Strategy
}

// 协程树根节点
type Supervisor struct {
	Name         string // supervisor的名称
	ex_channel   chan WorkerResponseMessage
	template     IWorker
	listen_mq    chan SupervisorReceiveMessage // supervisor监听该mq
	restart_rule Strategy
}

/*
supervisor只管理以下事务:
1. 开启goroutine
2. 关闭goroutine
3. 监听goroutine的失效情况
4. 在goroutine故障次数超过限制时，通知tracker，即主状态机
*/

func create_supervisor(option Option) chan SupervisorReceiveMessage {
	listen_mq := make(chan SupervisorReceiveMessage, 1000)
	monitor := &Supervisor{
		Name:         option.Name,
		tracker_mq:   option.Tracker_mq,
		template:     option.Template,
		listen_mq:    listen_mq,
		restart_rule: option.RestartRule,
	}
	go start_supervisor(monitor)
	return listen_mq
}

func start_supervisor(monitor *Supervisor) {
	// 构建一个dict，用于记录奔溃的goroutine的具体统计数据，比如时间和次数
	log.Println("Monitor start_supervisor start")
	// entrys := make(type EntryStatic struct {
	// 	Entry           string
	// 	Err_timestamps []int64
	// })
	entrys := make(map[string]EntryStatic)
	for {
		select {
		case message, ok := <-monitor.listen_mq:
			log.println("Monitor receive message")
			if !ok {
				log.Panic("Monitor receive invalid message")
			}
			if message.Act == CREATE_ROUTINE {
				entry, ok := entrys[message.Name]
				if !ok {
					worker_mq := make(chan WorkerReceiveMessage, 1000)
					go start_worker(monitor, worker_mq)
					entry[message.Name] = EntryStatic{
						Name: message.Name,
						Mq:   worker_mq,
					}

				}
				message.Mq <- entry.mq
			} else if message.Act == REMOVE_ROUTINE {

			} else if message.Act == DOWN_ROUTINE {

			} else {
				log.Panic("Monitor receive unknow message: ", message)
			}

			param := Registry.GetRegistryParameter(message)
			log.Println("Tracker start_supervisor receive message: ", message)
			if !ok {
				log.Println("Tracker start_supervisor receive channel close signal")
				return
			}
			_, ok = entry_static[param]
			if !ok {
				entry_static[param] = &Utils.EntryState{}
			}
			entry_static[param].Add()
			restart := track.Supervisor(entry_static[param])
			if restart {
				log.Println("Tracker start_supervisor restart success: ", message)
				track.RestartEntry(message)
			} else {
				// close the entry 以及对应的channel
				log.Println("Tracker start_supervisor restart fail: ", message)
				track.RemoveEntry(message)
			}
		}
	}
}
