package main

import (
	. "gopkg.in/metakeule/watcher.v1"
	"gopkg.in/metakeule/watcher.v1/notifier/stripgopath"
	"gopkg.in/metakeule/watcher.v1/notifier/zenity"
	. "gopkg.in/metakeule/watcher.v1/project"
)

func main() {
	compilers := CompilersForProject("/home/benny/Entwicklung/gopath/src/koelnart")
	projectWatcher := New(stripgopath.StripGoPath(zenity.Zenity()), compilers...)
	projectWatcher.Run()
	<-projectWatcher.Ready
	for {
	}
}
