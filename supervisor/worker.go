package supervisor

import "log"

// 协程使用的回调对象接口

type IWorker interface {
	// 用于构建新的IWorker对象，特别需要处理引用类型变量，否则导致严重错误
	Init(entry_name string, mq chan WorkerReceiveMessage) IWorker
	Get(status interface{}, mq chan WorkerReceiveMessage) (interface{}, error)
	Refresh(status interface{}, mq chan WorkerReceiveMessage) (IWorker, error)
}

// 另一种类型的IWorker可以为, 即完全由客户端处理即可
// type IWorker interface {
// 	Init(entry_name string, mq chan WorkerReceiveMessage) IWorker
// 	Get(status interface{}, mq chan WorkerReceiveMessage) (interface{}, error)
// 	Put(status interface{}, mq chan WorkerReceiveMessage) (IWorker, error)
// }

// type Supervisor struct {
// 	Name         string // supervisor的名称
// 	tracker_mq   chan WorkerResponseMessage
// 	worker     IWorker
// 	listen_mq    chan SupervisorReceiveMessage // supervisor监听该mq
// 	restart_rule Strategy
// }
// type Entry struct {
// 	e_mu    sync.Mutex
// 	e_close bool
// 	Name    string
// 	Mq      chan WorkerReceiveMessage
// }

// 执行自动更新的worker模型
func start_autoupdate_worker(monitor *Supervisor, entry *Entry, worker IWorker) {
	log.Println("Start worker with entry Name: ", entry.Name)
	status_machine := worker.Init(entry.Name, entry.Mq)
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
			} else if data.MessageType == WORKER_MESSAGE_NORMAL {
				// normal动作通常不改变status
				result, err := handle_normal(status_machine, entry.Mq, data.Data)
				if err != nil {
					is_open = false
					err_reason = EntryPanicError{}
					data.Mq <- err_reason
				} else {
					data.Mq <- result
				}
			} else if data.MessageType == WORKER_MESSAGE_REFRESH {
				// refresh 用于更新status
				new_status, err := handle_refresh(status_machine, entry.Mq, data)
				if err != nil {
					is_open = false
					err_reason = EntryPanicError{}
					data.Mq <- err_reason
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

func handle_normal(worker IWorker, entry_mq chan WorkerReceiveMessage, data interface{}) (result interface{}, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker Handle normal panic: ", r)
			err = EntryPanicError{}
		}
	}()
	result, err = worker.Get(data, entry_mq)
	return
}

func handle_refresh(worker IWorker, entry_mq chan WorkerReceiveMessage, data interface{}) (status IWorker, err error) {
	err = nil
	defer func() {
		if r := recover(); r != nil {
			log.Println("Worker Handle refresh panic: ", r)
			err = EntryPanicError{}
		}
	}()
	status, err = worker.Refresh(data, entry_mq)
	return
}
