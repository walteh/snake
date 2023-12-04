package swails

type WailsEvent struct {
	Name string
	Func func(args ...any)
}

func (me *WailsSnake) Events() []*WailsEvent {
	return []*WailsEvent{
		{
			Name: "wails:ready",
			Func: func(args ...any) {

			},
		},
	}
}
