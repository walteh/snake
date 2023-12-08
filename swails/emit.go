package swails

import (
	"context"
)

// WailsEmitter is the function type as provided
type WailsEmitter func(ctx context.Context, eventName string, optionalData ...interface{})

type WailsEmitEvent struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

// WailsWriter implements the io.Writer interface and sends data to a WailsEmitter.
type WailsWriter struct {
	emitter WailsEmitter
	Name    string `json:"name"`
	context context.Context
}

func (me *WailsSnake) noopWailsEmitEvent() *WailsEmitEvent {
	return &WailsEmitEvent{}
}

// Write implements the io.Writer interface for WailsWriter.
func (w *WailsWriter) Write(p []byte) (n int, err error) {
	// Convert the byte slice to a string or any desired format
	data := string(p)

	dat := &WailsEmitEvent{
		Data: data,
		Name: w.Name,
	}

	// dato, err := json.Marshal(dat)
	// if err != nil {
	// 	return 0, err
	// }

	// data = string(dato)

	// fmt.Printf("Writing data to %q: %q\n", w.Name, data)

	// Emit the data using the WailsEmitter function
	w.emitter(w.context, w.Name, dat)

	// Return the number of bytes written and no error
	return len(p), nil
}

func (me *WailsSnake) newEventEmitter(id string) *WailsWriter {
	return &WailsWriter{
		emitter: me.emitter,
		Name:    id,
		context: me.lifecycleContext,
	}
}
