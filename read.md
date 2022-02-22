通过该库构建goroutine tree，该库不但构建entry对象**核心为一个channel和goroutine对**，并传递给外部调用方
同时监控它们健康状态，按照给定的规则决定是否重启以及报告

提供接口
1. AddWorker, return channel
2. RemoveWorker， return bool

目前chan的长度固定，后续可以改为可配置

go_registry实现中，使用了supervisor或者worker通知tracker去除entry项的操作
但是，此处实现不准备处理该步骤，直接由用户调用的时候，通过entry的close选项确定


