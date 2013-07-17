package zenity

import (
	"fmt"
	. "github.com/metakeule/watcher/helpers"
	"log"
)

type zenityBin string

func (ø zenityBin) Error(msg string) {
	Exec(string(ø), "--info", fmt.Sprintf("--text=%s", msg))
}

func (ø zenityBin) Success(msg string) {
	//Exec(string(ø), "--info", fmt.Sprintf("--text=%s", msg))
}

func New() zenityBin {
	path, err := Which("zenity")
	if err != nil {
		log.Fatalf("can't find zenity")
	}
	return zenityBin(path)
}
