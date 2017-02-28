package protocol

const (
    _  = 0x00
    Register = iota
    RegisterSuccess
    RegisterFailed

    SubmitTask
    SubmitTaskSuccess
    SubmitTaskFailed

    Ping
    Pong

    Trace
    TraceSuccess
    TraceFailed

    StopTask
    StopTaskSuccess
    StopTaskFailed
)