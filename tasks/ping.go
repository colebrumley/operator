package tasks

// Ping is a super basic task that returns "Pong"
func Ping(...string) (string, error) {
	return "Pong", nil
}
