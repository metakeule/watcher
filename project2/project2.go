package project2

import (
	. "gopkg.in/metakeule/watcher.v1"
	"gopkg.in/metakeule/watcher.v1/app2"
	. "gopkg.in/metakeule/watcher.v1/compiler/goc"
	"io/ioutil"
	"log"
	"path"
)

func compilersForAllApps(baseDir string) (cs []Compiler) {
	dirs, e := ioutil.ReadDir(path.Join(baseDir, "app"))
	if e != nil {
		log.Fatalf("could not read dir %s: %s", baseDir, e.Error())
	}
	cs = []Compiler{}
	for _, dir := range dirs {
		if dir.IsDir() == true {
			cs = append(cs, app2.Compilers(baseDir, dir.Name())...)
		}
	}
	return
}

func Compilers(execTests bool, baseDir string) (compilers []Compiler) {
	runDir := path.Join(baseDir, "run")
	mainFile := path.Join(runDir, "main.go")
	// pidFile := path.Join(runDir, "pid")
	compilers = []Compiler{
		NewGoCompiler(mainFile, baseDir, execTests, []string{"static", "less", "typescript"}),
	}
	compilers = append(compilers, compilersForAllApps(baseDir)...)
	return
}
