package worker

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}
