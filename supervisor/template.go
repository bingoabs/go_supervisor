package supervisor

// 协程使用的回调对象接口

type IWorker interface {
	// 定义初始化状态, 在注册goroutine时执行, 若定义定时任务, 那么需要Handler中进行处理REFRESH信息
	Init(message InitMessage) interface{}
	// 执行goroutine接收到的请求，返回值为 {新的status, 返回值内容(如果不需要返回值，需为nil)}
	// 注意, 实际上即便是异步请求也可以返回内容, 由用户自己决定
	// 以及定时更新，可以通过该函数进行处理，即对该函数暴露了该协程对应的通道
	Handler(message RequestMessage, mq chan RequestMessage) (interface{}, interface{})
}
