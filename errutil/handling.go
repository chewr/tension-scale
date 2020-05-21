package errutil

func SwallowF(f func() error) {
	_ = f()
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
