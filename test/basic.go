package test

import (
	Supervisor "github.com/bingoabs/go_supervisor/supervisor"
)

func getSupervisorOption() Supervisor.Option {
	option := Supervisor.Option{
		Name:      "test_supervisor",
		Generator: CreateHttpWorker,
		RestartRule: Supervisor.Strategy{
			TimeInterval: 10,
			Number:       3,
		},
		RefreshInterval: 5,
	}
	return option
}

func getPanicWorkerSupervisorOption() Supervisor.Option {
	option := Supervisor.Option{
		Name:      "test_supervisor",
		Generator: CreateRefreshPanicHttpWorker,
		RestartRule: Supervisor.Strategy{
			TimeInterval: 15,
			Number:       3,
		},
		RefreshInterval: 2,
	}
	return option
}
