package supervisor

// 要求在给定事件间隔期间，不能失败超过Number次
type Strategy struct {
	TimeInterval int64
	Number       int
}
