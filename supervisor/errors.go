package supervisor

type ClosedWorkerError struct {
}

func (e ClosedWorkerError) Error() string {
	return "Worker already closed"
}

type EntryPanicError struct {
}

func (e EntryPanicError) Error() string {
	return "Entry goroutine already panic"
}
