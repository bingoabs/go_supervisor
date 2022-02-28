通过该库构建goroutine tree，该库不但构建entry对象**核心为一个channel和goroutine对**，并传递给外部调用方
同时监控它们健康状态，按照给定的规则决定是否重启以及报告

提供接口
1. AddWorker, return channel
2. RemoveWorker， return bool

目前chan的长度固定，后续可以改为可配置

go_registry实现中，使用了supervisor或者worker通知tracker去除entry项的操作
但是，此处实现不准备处理该步骤，直接由用户调用的时候，通过entry的close选项确定


总的来说，go_supervisor是一个goroutine管理器，用于创建满足特定接口的对象作为内部状态的无限循环监听特定channle的goroutine
外部调用对象得到的是一个entry对象，该对象包含了goroutine的channel和状态信息，并通过entry与goroutine通信
在goroutine失效后，go_supervisor会通过channel通知监听该事件的外部对象，注意，如果goroutine失效的情况下满足重启条件，那么go_supervisor会重启goroutine，并提供了幂等的create接口


// 用户通过对实现的IWorker对象设置字段, 从而实现IWorker对象本身为一个状态机
// 每个goroutine维护一个IWorker,因此 IWorker的字段不能存在指针
// 如果一定要使用, 那么在Init函数中必须处理好


worker负责记录当前goroutine的失败次数，并检查是否超过限制，如果超过，那么改变goroutine状态，并开始返回错误信息，而不是直接关闭channel，并通知supervisor
而supervisor通过监听down信息，清除entry信息，并执行entry的关闭

注意，supervisor、
go_registry 拆分为 go_registry, go_supervisor, go_mesh, 其中三者可以成为goroutine注册中心，而go_supervisor单独作为协程管理组件，go_mesh单独作为节点间通信组件
supervisor只是管理goroutine，以及在panic超过容忍限制后崩溃，并不关心其他

而 最终的 go_registy实现的是在各个节点分布goroutine，但是每个goroutine必须保证可以直接替换，即goroutine本身逻辑保证内容，而不能依赖节点间的重排列过程会保持goroutine 的状态完美切换

也就是说，用于实现缓存很好用，但是如果做数据持久化可能不是什么好主意
						// 该work用于自动从远端更新内容,还可以执行其他实现,比如由用户对status进行更新

