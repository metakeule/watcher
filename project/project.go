package project

import (
	"io/ioutil"
	"log"
	"path"

	"gopkg.in/metakeule/watcher.v1"
	"gopkg.in/metakeule/watcher.v1/app"
	"gopkg.in/metakeule/watcher.v1/compiler/goc"
)

func compilersForAllApps(baseDir string) (cs []watcher.Compiler) {
	dirs, e := ioutil.ReadDir(path.Join(baseDir, "app"))
	if e != nil {
		log.Fatalf("could not read dir %s: %s", baseDir, e.Error())
	}
	cs = []watcher.Compiler{}
	for _, dir := range dirs {
		if dir.IsDir() == true {
			cs = append(cs, app.Compilers(baseDir, dir.Name())...)
		}
	}
	return
}

func Compilers(execTests bool, baseDir string) (compilers []watcher.Compiler) {
	runDir := path.Join(baseDir, "run")
	mainFile := path.Join(runDir, "main.go")
	// pidFile := path.Join(runDir, "pid")
	compilers = []watcher.Compiler{
		goc.NewGoCompiler(mainFile, baseDir, execTests, []string{"static", "less", "typescript"}),
	}
	compilers = append(compilers, compilersForAllApps(baseDir)...)
	return
}
