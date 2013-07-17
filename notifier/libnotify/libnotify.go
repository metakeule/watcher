package libnotify

import (
	. "github.com/metakeule/watcher/helpers"
	"log"
)

type notifyBin string

func (ø notifyBin) Error(msg string) {
	Exec(string(ø), "ERROR", msg)
}

func (ø notifyBin) Success(msg string) {
	Exec(string(ø), "Ok", msg)
}

func New() notifyBin {
	path, err := Which("notify-send")
	if err != nil {
		log.Fatalf("can't find notify-send, under ubuntu run sudo apt-get install libnotify-bin")
	}
	return notifyBin(path)
}
