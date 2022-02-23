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


go中不存之前缀++--，只能使用后缀，并且不是表达式，即 a = j++是非法的

// TODO, 关闭的channel是否能读尽内容 可以

Implementing interface in golang gives method has pointer receiver [duplicate]
指针类型实现的接口，那么就是指针类型可以作为该接口类型使用，但是指针所指向的类型不行！


// // empty slice 是指slice不为nil，但是slice没有值，slice的底层的空间是空
// slice := make([]int, 0) // slice := []int{}
// var slice []int wei nil
// golang中允许对值为 nil 的 slice 添加元素

// myMap = make(map[string] personInfo, 5)
panic: assignment to entry in nil map