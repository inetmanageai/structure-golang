package events

// NOTE PORT -----------------------------------------------
type AppEvent interface {
	On(topics string, group string, handler func(topic string, message []byte) error) (err error)

	Emit(topic string, data []byte) (patition int32, offset int64, err error)
}
