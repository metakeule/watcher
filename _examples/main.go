package main

import (
	. "github.com/metakeule/watcher"
	"github.com/metakeule/watcher/notifier/stripgopath"
	"github.com/metakeule/watcher/notifier/zenity"
	. "github.com/metakeule/watcher/project"
)

func main() {
	compilers := CompilersForProject("/home/benny/Entwicklung/gopath/src/koelnart")
	projectWatcher := New(stripgopath.StripGoPath(zenity.Zenity()), compilers...)
	projectWatcher.Run()
	<-projectWatcher.Ready
	for {
	}
}
