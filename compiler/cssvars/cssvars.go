package cssvars

import (
	"github.com/metakeule/goh4"
	. "github.com/metakeule/watcher/helpers"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type cssVarsCompiler struct {
	*sync.Mutex
	All        []goh4.Class
	Package    string
	Dir        string
	OutputFile string
	Bin        string
	name       string
	Var        string
}

func (ø *cssVarsCompiler) Name() string              { return "cssVars compiler (" + ø.name + ")" }
func (ø *cssVarsCompiler) Dirs() []string            { return []string{ø.Dir} }
func (ø *cssVarsCompiler) Affected(path string) bool { return IsInDir(path, ø.Dir) }
func (ø *cssVarsCompiler) Compile(path string) (output string, err error) {
	// -in=koelnart/frontend/style/class -target=/home/benny/Entwicklung/gopath/src/koelnart/frontend/less/classes.less
	if ø.Var == "" {
		output, err = Exec(ø.Bin, "-in="+ø.Package, "-target="+ø.OutputFile)
	} else {
		output, err = Exec(ø.Bin, "-in="+ø.Package, "-target="+ø.OutputFile, "-var="+ø.Var)
	}
	return
}

func New(name string, dir string, output string, var_ string) (ø *cssVarsCompiler) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("GOPATH is not set")
	}

	var err error
	ø = &cssVarsCompiler{}
	ø.name = name
	ø.Bin, err = Which("cssvars")
	if err != nil {
		log.Fatal("you do not have cssvars installed. run go get github.com/metakeule/goh4/css/cssvars")
	}

	/*
	   if !IsInDir(gopath, dir) {
	       log.Fatalf("%s is not in $GOPATH", dir)
	   }
	*/
	d, e := os.Stat(dir)
	if e != nil {
		log.Fatalf("%s does not exists: %s", dir, e.Error())
	}
	if !d.IsDir() {
		log.Fatalf("%s is no package", dir)
	}
	ø.Dir = dir
	packageName := strings.Replace(dir, path.Join(gopath, "src")+"/", "", 1)
	ø.Package = packageName
	ø.Var = var_
	ø.OutputFile = output
	ø.Mutex = &sync.Mutex{}
	return
}
