package internal

var SchedulerDefaultArgs = []string{
	"--master={{ .APIServerURL.String }}",
	"--port={{ .URL.Port }}",
	"--address={{ .URL.Hostname }}",
}

func DoSchedulerArgDefaulting(args []string) []string {
	if len(args) != 0 {
		return args
	}

	return SchedulerDefaultArgs
}
