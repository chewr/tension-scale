package display

type fnStateE func(State) error

func (fn fnStateE) UpdateState(state State) error {
	return fn(state)
}

// TODO(rchew) States are mutable now, so muxing is weird... A single StateHolder that is shared between multiple displays seems better
// -> I need to refactor displays to take a StateHolder, then delete this file
func ModelMux(models ...Model) Model {
	return fnStateE(func(state State) error {
		for _, model := range models {
			if err := model.UpdateState(state); err != nil {
				return err
			}
		}
		return nil
	})
}
