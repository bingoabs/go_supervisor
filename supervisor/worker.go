package supervisor

import "log"

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
	log.Println("Start worker with entry Name: ", entry.Name)
	status_machine := worker.Init(entry)
	is_open := true
	var err_reason error
	for {
		select {
		// 此处不能使用context直接结束无限循环，因为Mq中可能存在未处理消息，必须全部处理完才行
		case data, ok := <-entry.Mq:
			log.Println("Worker receive message")
			if !ok {
				log.Println("Worker receive channel close signal, break")
				return
			}
			if !is_open {
				// 用户在接收到错误信息后，应该弃用当前持有entry，并从tracker获取新的entry
				// 继续使用旧的entry, 后果自负
				data.Mq <- err_reason
			} else if data.MessageType == WORKER_MESSAGE_STOP {
				is_open = false
				err_reason = ClosedWorkerError{}
			} else if data.MessageType == WORKER_MESSAGE_GET {
				result, err := handle_get(status_machine, data.Data)
				if err != nil {
					is_open = false
					err_reason = EntryPanicError{}
					data.Mq <- err_reason
				} else {
					data.Mq <- result
				}
			} else if data.MessageType == WORKER_MESSAGE_PUT {
				new_status, err := handle_put(status_machine, data)
				if err != nil {
					is_open = false
					err_reason = EntryPanicError{}
					data.Mq <- err_reason
				} else {
					status_machine = new_status
					data.Mq <- true
				}
			} else if data.MessageType == WORKER_MESSAGE_REFRESH {
				// refresh 用于更新status
				new_status, err := handle_refresh(status_machine, entry, data, monitor.refresh_interval)
				if err != nil {
					is_open = false
					err_reason = EntryPanicError{}
				} else {
					status_machine = new_status
				}
			} else {
				log.Panic("Worker receive unkonwn message type: ", data.MessageType)
			}
		}
	}
}

func close_worker(monitor *Supervisor, entry *Entry) {

}

func handle_get(worker IWorker, data interface{}) (result interface{}, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker Handle get panic: ", r)
			err = EntryPanicError{}
		}
	}()
	result, err = worker.Get(data)
	return
}
func handle_put(worker IWorker, data interface{}) (status IWorker, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker Handle put panic: ", r)
			err = EntryPanicError{}
		}
	}()
	status, err = worker.Put(data)
	return
}

func handle_refresh(worker IWorker, entry *Entry, data interface{}, interval int) (status IWorker, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker Handle refresh panic: ", r)
			err = EntryPanicError{}
		}
	}()
	status, err = worker.Refresh(data)
	if interval > 0 {
		go send_refersh_signal(entry, interval)
	}
	return
}
