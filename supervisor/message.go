package supervisor

// 工作goroutine发送给调用方的消息格式
type WorkerResponseMessage struct {
	Err  error
	Data interface{}
}

// worker通道相关信息, 类型为私有，但是设定常量为公开
type workerMessageType string

// 实际上无法实现真正的异步，因为即便是异步操作也需要确定该goroutine是否仍然工作，因此全部同步
// 如果有异步的操作，需要用户在自定义worker中实现相关逻辑
const WORKER_MESSAGE_STOP workerMessageType = "STOP"
const WORKER_MESSAGE_GET workerMessageType = "GET"
const WORKER_MESSAGE_REFRESH workerMessageType = "REFRESH"

// 工作goroutine接收消息的格式, 目前直接使用 interface{}
type WorkerReceiveMessage struct {
	MessageType workerMessageType
	Data        interface{}
	Mq          chan WorkerResponseMessage
}

// worker和supervisor监听的channel的长度设置值
const WORKER_MQ_LENGTH int = 1000
const SUPERVISOR_MQ_LENGTH int = 1000

// Supervisor通道相关信息
type supervisorMessageType string

const (
	SUPERVISOR_CREATE_EVENT supervisorMessageType = "CREATE_EVENT"
	SUPERVISOR_REMOVE_EVENT supervisorMessageType = "REMOVE_EVENT"
	SUPERVISOR_DOWN_EVENT   supervisorMessageType = "DOWN_EVENT"
)

type SupervisorReceiveMessage struct {
	EntryName   string
	MessageType supervisorMessageType
	Mq          chan *Entry
}
