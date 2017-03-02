package protocol

const (
	Join = iota + 1
	JoinAccept
	JoinReject

	TaskSubmit
	TaskSumbitAccept
	TaskSubmitReject

	Ping
	Pong

	TaskStop
	TaskStopAccept
	TaskStopReject
)
