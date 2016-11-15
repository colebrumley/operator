package tasks

// Task represents a single action that Operator can perform.
type Task struct {
	Description string
	Fn          func(...string) (string, error)
}

// TaskList is a map of available task handlers. It is used to dynamically determine API routes.
var TaskList = map[string]Task{
	"ping": Task{
		Description: "A simple task that returns 'Pong'",
		Fn:          Ping,
	},
	"ssh": Task{
		Description: "Dispatch a command to a remote node via SSH",
		Fn:          SSH,
	},
}
