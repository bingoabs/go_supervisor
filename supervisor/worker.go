package supervisor

import "log"

// 协程使用的回调对象接口

// 用户通过对实现的IWorker对象设置字段, 从而实现IWorker对象本身为一个状态机
// 每个goroutine维护一个IWorker,因此 IWorker的字段不能存在指针
// 如果一定要使用, 那么在Init函数中必须处理好

// todo, 使用IWorker替换status
// TODO, 关闭的channel是否能读尽内容
// TODO, 接口的指针是否可以调用接口函数
type IWorker interface {
	Init(entry_name string, mq chan WorkerReceiveMessage) IWorker
	Get(status interface{}, mq chan WorkerReceiveMessage) (interface{}, error)
	Update(status interface{}, mq chan WorkerReceiveMessage) (interface{})
}

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
	status_machine := worker.Init(entry.Name)
	is_open := true
	var err_reason error
	for {
		select {
		// 此处不能使用context直接结束无限循环，因为Mq中可能存在未处理消息，必须全部处理完才行
		case data, ok := <-entry.Mq:
			log.Println("Worker receive message")
			if !ok {
				log.Println("Worker receive channel close signal, break")
				break
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
				new_status, err := handle_refresh(worker, entry.Mq, status, data)
				if err != nil {
					is_open = false 
					err_reason = EntryPanicError{}
					data.Mq <- err_reason
				} else {
					status = new_status
				}
			}else {
				log.Panic("Worker receive unkonwn message type: ", data.MessageType)
			}
		}
	}
}

func close_worker(monitor *Supervisor, entry *Entry) {
}

func handle_normal(worker IWorker, entry_mq chan WorkerReceiveMessage, status, data interface{}) (interface{}, error) {
	
}

func handle_refresh(worker IWorker, entry_mq chan WorkerReceiveMessage, status, data interface{}) (interface{}, error) {
	
}

func handle_worker_request(worker IWorker, status interface{}, data interface{}, ) (status interface{}, interface{}) {
	new_status
	return nil
}
