package libnotify

import (
	. "gopkg.in/metakeule/watcher.v1/helpers"
	"log"
)

type notifyBin string

func (ø notifyBin) Error(msg string) {
	log.Printf("ERROR: %s", msg)
	Exec(string(ø), "--icon=dialog-error", "--expire-time=1004000", "ERROR", msg)
}

func (ø notifyBin) Success(msg string) {
	//Exec(string(ø), "--icon=dialog-information", "--expire-time=2500", "Ok", msg)
	//Exec(string(ø), "--expire-time=2500", "Ok", msg)
	Exec(string(ø), "--expire-time=2000", msg)
}

func New() notifyBin {
	path, err := Which("notify-send")
	if err != nil {
		log.Fatalf("can't find notify-send, under ubuntu run sudo apt-get install libnotify-bin")
	}
	return notifyBin(path)
}
