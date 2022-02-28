package supervisor

// 协程使用的回调对象接口

// 如果设置更新时间间隔为0，那么该IWorker退化为简单的状态机
type IWorker interface {
	Get(status interface{}) (interface{}, error)
	Put(status interface{}) (IWorker, error)
	Refresh(status interface{}) (IWorker, error)
}

// 用于构建新的IWorker对象，特别需要处理引用类型变量，否则导致严重错误
type WorkerGenerator func(entry *Entry) IWorker

// supervisor调用接口
func CreateSupervisor(option Option) chan SupervisorReceiveMessage {
	return create_supervisor(option)
}

func GetWorker(mq chan SupervisorReceiveMessage, entry_name string) *Entry {
	message := SupervisorReceiveMessage{
		EntryName:   entry_name,
		MessageType: SUPERVISOR_CREATE_EVENT,
		Mq:          make(chan *Entry, 1),
	}
	return send_message_to_supervisor(mq, message)
}

func RemoveWorker(mq chan SupervisorReceiveMessage, entry_name string) *Entry {
	message := SupervisorReceiveMessage{
		EntryName:   entry_name,
		MessageType: SUPERVISOR_REMOVE_EVENT,
		Mq:          make(chan *Entry, 1),
	}
	return send_message_to_supervisor(mq, message)
}
