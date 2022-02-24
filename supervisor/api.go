package supervisor

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
