package display

type State interface{}

type Model interface {
	UpdateState(state State)
}
