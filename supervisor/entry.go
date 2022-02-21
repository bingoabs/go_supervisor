package supervisor

type Entry struct {
	Name string
	Mq   chan WorkerReceiveMessage
}

type EntryStatic struct {
	Entry          *Entry
	Err_timestamps []int64
}
