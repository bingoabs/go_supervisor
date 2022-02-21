package supervisor

// supervisor调用接口
func Create(option Option) chan SupervisorReceiveMessage {
	return create_supervisor(option)
}
func Insert(supervisor_mq chan SupervisorReceiveMessage) {

}
