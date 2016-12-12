package tasks

// Ping is a super basic task that returns "Pong"
func Ping(...interface{}) (string, error) {
	return "Pong", nil
}
