package internal

var ControllerManagerDefaultArgs = []string{
	"--master={{ .APIServerURL.String }}",
	"--port={{ .URL.Port }}",
	"--address={{ .URL.Hostname }}",
}

func DoControllerManagerArgDefaulting(args []string) []string {
	if len(args) != 0 {
		return args
	}

	return ControllerManagerDefaultArgs
}
