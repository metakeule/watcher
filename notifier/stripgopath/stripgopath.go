package stripgopath

import (
	"github.com/metakeule/watcher"
	"log"
	"os"
	"path"
	"strings"
)

type stripGoPath struct {
	watcher.Notifier
	GoPath string
}

func New(n watcher.Notifier) watcher.Notifier {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("GOPATH is not set")
	}
	return &stripGoPath{n, path.Join(gopath, "src") + "/"}
}

func (ø *stripGoPath) Notify(msg string) {
	nm := strings.Replace(msg, ø.GoPath, "", -1)
	ø.Notifier.Notify(nm)
}
