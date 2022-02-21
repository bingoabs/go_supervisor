package supervisor

// 工作goroutine发送给调用方的消息格式
type WorkerResponseMessage struct {
	Err  error
	Data interface{}
}

// 工作goroutine接收消息的格式
type WorkerReceiveMessage struct {
	Data interface{}
}

// 为了限制外部使用不知名的action，设置如下
type Action uint32

const (
	CREATE_ROUTINE Action = iota
	UPDATE_ROUTINE
	REMOVE_ROUTINE
	DOWN_ROUTINE
)

type SupervisorReceiveMessage struct {
	Name string
	Act  Action
	Mq   chan interface{}
}
