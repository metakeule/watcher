package project

import (
	. "github.com/metakeule/watcher"
	"github.com/metakeule/watcher/app"
	. "github.com/metakeule/watcher/compiler/goc"
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
			cs = append(cs, app.Compilers(baseDir, dir.Name())...)
		}
	}
	return
}

func Compilers(execTests bool, baseDir string) (compilers []Compiler) {
	mainFile := path.Join(baseDir, "main.go")
	pidFile := path.Join(baseDir, "pid")
	compilers = []Compiler{
		NewGoCompiler(mainFile, pidFile, baseDir, execTests, []string{"static", "less", "typescript"}),
	}
	compilers = append(compilers, compilersForAllApps(baseDir)...)
	return
}
