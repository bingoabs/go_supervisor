package supervisor

import (
	"log"
	"time"
)

// 协程使用的回调对象接口

// 如果设置更新时间间隔为0，那么该IWorker退化为简单的状态机
type IWorker interface {
	// 用于构建新的IWorker对象，特别需要处理引用类型变量，否则导致严重错误
	Init(entry *Entry) IWorker
	Get(status interface{}) (interface{}, error)
	Put(status interface{}) (IWorker, error)
	Refresh(status interface{}) (IWorker, error)
}

// 执行自动更新的worker模型
func start_autoupdate_worker(monitor *Supervisor, entry *Entry, worker IWorker) {
	log.Println("Worker start_autoupdate_worker start with entry Name: ", entry.Name)
	status_machine := worker.Init(entry)
	is_open := true
	var err_reason error
	panic_timestamps := []int64{}
	if monitor.refresh_interval > 0 {
		go refresh_worker(entry, monitor.refresh_interval)
	}
	for {
		// 此处不能使用context直接结束无限循环，因为Mq中可能存在未处理消息，必须全部处理完才行
		// 因此仅一个分支，不需要使用select
		data, ok := <-entry.Mq
		log.Println("Worker start_autoupdate_worker receive message")
		if !ok {
			log.Println("Worker start_autoupdate_worker receive close signal and break")
			return
		}
		// 用户在接收到错误信息后，应该弃用当前持有entry，并从tracker获取新的entry; 继续使用旧的entry, 后果自负
		if !is_open {
			log.Println("Worker start_autoupdate_worker already closed")
			data.Mq <- WorkerResponseMessage{
				Err:  err_reason,
				Data: nil,
			}
		} else if data.MessageType == WORKER_MESSAGE_STOP {
			log.Println("Worker start_autoupdate_worker receive STOP event")
			is_open = false
			err_reason = ClosedWorkerError{}
			data.Mq <- WorkerResponseMessage{
				Err:  nil,
				Data: true,
			}
		} else if data.MessageType == WORKER_MESSAGE_GET {
			log.Println("Worker start_autoupdate_worker receive GET event")
			// 直接捕捉panic, 不必通过supervisor的channel消息兜一圈
			result, err := handle_get(status_machine, data.Data)
			if err != nil {
				log.Println("Worker start_autoupdate_worker GET event receive err: ", err)

				panic_timestamps = append(panic_timestamps, time.Now().Unix())
				if valid_panic_times(monitor.restart_rule, panic_timestamps) {
					data.Mq <- WorkerResponseMessage{
						Err:  err,
						Data: nil,
					}
				} else {
					is_open = false
					err_reason = err
					data.Mq <- WorkerResponseMessage{
						Err:  err_reason,
						Data: nil,
					}
					go send_down_to_supervisor(monitor, entry)
				}
			} else {
				log.Println("Worker start_autoupdate_worker GET event receive result: ", result)
				data.Mq <- WorkerResponseMessage{
					Err:  nil,
					Data: result,
				}
			}
		} else if data.MessageType == WORKER_MESSAGE_REFRESH {
			log.Println("Worker start_autoupdate_worker receive REFRESH event")

			// refresh 用于更新status
			result, err := handle_refresh(status_machine, entry, data, monitor.refresh_interval)
			if err != nil {
				log.Println("Worker start_autoupdate_worker REFRESH event receive err: ", err)

				panic_timestamps = append(panic_timestamps, time.Now().Unix())
				if valid_panic_times(monitor.restart_rule, panic_timestamps) {
					data.Mq <- WorkerResponseMessage{
						Err:  err,
						Data: nil,
					}
				} else {
					is_open = false
					err_reason = err
					data.Mq <- WorkerResponseMessage{
						Err:  err_reason,
						Data: nil,
					}
					go send_down_to_supervisor(monitor, entry)
				}
			} else {
				log.Println("Worker start_autoupdate_worker REFRESH event receive result: ", result)
				status_machine = result
				data.Mq <- WorkerResponseMessage{
					Err:  nil,
					Data: true,
				}
			}
		} else {
			log.Panic("Worker receive unkonwn message type: ", data.MessageType)
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
func close_worker(monitor *Supervisor, entry *Entry) {
	log.Println("Worker close_worker start")
	message := WorkerReceiveMessage{
		MessageType: WORKER_MESSAGE_STOP,
		Data:        nil,
		Mq:          make(chan WorkerResponseMessage, 1),
	}
	result := send_message_to_worker(entry, message)
	log.Println("Worker close_worker receive: ", result)
}

func handle_get(worker IWorker, data interface{}) (result interface{}, err error) {
	log.Println("Worker handle_get start")
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker Handle get panic: ", r)
			err = EntryPanicError{}
		}
	}()
	result, err = worker.Get(data)
	log.Println("Worker handle_get result: ", result)
	return
}

func handle_refresh(worker IWorker, entry *Entry, data interface{}, interval int) (status IWorker, err error) {
	log.Println("Worker handle_refresh start")
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker handle_refresh panic reason: ", r)
			err = EntryPanicError{}
		}
	}()
	status, err = worker.Refresh(data)
	if interval > 0 {
		go refresh_worker(entry, interval)
	}
	log.Println("Worker handle_refresh end")
	return
}

func refresh_worker(entry *Entry, interval int) {
	log.Println("Worker refresh_worker start with interval: ", interval)
	ticker := time.NewTimer(time.Second * time.Duration(interval))
	<-ticker.C
	message := WorkerReceiveMessage{
		MessageType: WORKER_MESSAGE_REFRESH,
		Data:        nil,
		Mq:          make(chan WorkerResponseMessage, 1),
	}
	result := send_message_to_worker(entry, message)
	log.Println("Worker refresh_worker receive result: ", result)
}

func send_down_to_supervisor(monitor *Supervisor, entry *Entry) {
	message := SupervisorReceiveMessage{
		EntryName:   entry.Name,
		MessageType: SUPERVISOR_DOWN_EVENT,
		Mq:          make(chan *Entry, 1),
	}
	result := send_message_to_supervisor(monitor.listen_mq, message)
	log.Println("Worker send_down_to_supervisor result: ", result)
}
