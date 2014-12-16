package zenity

import (
	"fmt"
	. "gopkg.in/metakeule/watcher.v1/helpers"
	"log"
)

type zenityBin string

func (ø zenityBin) Error(msg string) {
	log.Printf("ERROR: %s", msg)
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
