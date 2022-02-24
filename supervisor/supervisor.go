package supervisor

import (
	"log"
	"time"
)

type Option struct {
	Name            string
	Worker          IWorker
	RestartRule     Strategy
	RefreshInterval int
}

// 协程树根节点
type Supervisor struct {
	Name             string                        // supervisor的名称
	worker           IWorker                       // worker的具体执行对象
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
		worker:           option.Worker,
		listen_mq:        listen_mq,
		restart_rule:     option.RestartRule,
		refresh_interval: option.RefreshInterval,
	}
	go start_supervisor(monitor)
	return listen_mq
}

func start_supervisor(monitor *Supervisor) {
	// 构建一个dict，用于记录奔溃的goroutine的具体统计数据，比如时间和次数
	log.Println("Monitor start_supervisor start")
	entry_statics := make(map[string]*EntryStatic)
	for {
		select {
		case message, ok := <-monitor.listen_mq:
			log.Println("Monitor receive message")
			if !ok {
				log.Panic("Monitor receive invalid message")
				return
			}
			// create, remove 由tracker调用, 而panic由定制worker调用
			if message.MessageType == CREATE_ROUTINE {
				log.Println("Monitor receive CREATE ROUTINE action")

				_, ok := entry_statics[message.EntryName]
				if !ok {
					worker_mq := make(chan WorkerReceiveMessage, WORKER_MQ_LENGTH)
					entry := &Entry{
						Name:       message.EntryName,
						Mq:         worker_mq,
						Created_At: time.Now().Unix(),
						e_closed:   false,
					}
					entry_statics[message.EntryName] = &EntryStatic{
						Entry: entry,
					}
					go start_autoupdate_worker(monitor, entry, monitor.worker)
				}
				message.Mq <- entry_statics[message.EntryName].Entry
			} else if message.MessageType == REMOVE_ROUTINE {
				log.Println("Monitor receive REMOVE ROUTINE action")

				_, ok := entry_statics[message.EntryName]
				if ok {
					entry := entry_statics[message.EntryName].Entry
					delete(entry_statics, message.EntryName)
					go close_worker(monitor, entry)
				}
			} else if message.MessageType == DOWN_ROUTINE {
				log.Println("Monitor receive DOWN ROUTINE action")

				_, ok := entry_statics[message.EntryName]
				if ok {
					// 每一次崩溃都被记录，如果崩溃次数超过限制，那么执行终结该worker, 否则重启
					entry_static := entry_statics[message.EntryName]
					entry_static.Panic_timestamps = append(entry_static.Panic_timestamps, time.Now().Unix())

					entry := entry_statics[message.EntryName].Entry
					go panic_worker(monitor, entry)
					// 如果down次数在合理区间，执行重建；否则删除记录
					// 特别注意，不再主动通知tracker，而是由调用Entry时进行检查，减少复杂性
					if valid_panic_times(monitor.restart_rule, entry_static.Panic_timestamps) {
						log.Println("Monitor receive valid DOWN event")
						worker_mq := make(chan WorkerReceiveMessage, WORKER_MQ_LENGTH)
						entry := &Entry{
							Name:     message.EntryName,
							Mq:       worker_mq,
							e_closed: false,
						}
						entry_statics[message.EntryName] = &EntryStatic{
							Entry:            entry,
							Panic_timestamps: entry_static.Panic_timestamps,
						}

						// TODO 使用配置选择不同的worker模型
						go start_autoupdate_worker(monitor, entry, monitor.worker)
					} else {
						log.Println("Monitor receive invalid DOWN event")
						delete(entry_statics, message.EntryName)
					}
				}
			} else {
				log.Panic("Monitor receive unknow message: ", message)
			}
		}
	}
}

func valid_panic_times(rule Strategy, timestamps []int64) bool {
	panics := len(timestamps)
	if rule.Number == 0 {
		return false
	}
	if panics <= rule.Number {
		return true
	}
	if timestamps[panics-1]-timestamps[panics-rule.Number] < int64(rule.TimeInterval) {
		return true
	}
	return false
}
