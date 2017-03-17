package scheduler

func Run() {
	var sched = newDkvScheduler()
	sched.Main()
}
