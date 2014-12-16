package wrapper

import (
	"gopkg.in/metakeule/watcher.v1"
)

type wrapper struct {
	// the inner notifier
	watcher.Notifier

	// a function that receives the message and a boolean, indicating, if
	// compilation was a success
	// it returns a boolean that should be true, if the inner Notifier
	// should be called and otherwise false
	fn func(msg string, success bool) (shouldContinue bool)
}

func New(n watcher.Notifier, fn func(msg string, success bool) (shouldContinue bool)) watcher.Notifier {
	return &wrapper{n, fn}
}

func (ø *wrapper) Error(msg string) {
	if ø.fn(msg, false) {
		ø.Notifier.Error(msg)
	}
}

func (ø *wrapper) Success(msg string) {
	if ø.fn(msg, true) {
		ø.Notifier.Success(msg)
	}
}
