package frames

// Frame хранит кадр, состоящий из слайса строк.
type Frame []string

func NewFrame(length int) Frame {
	return make(Frame, length)
}
