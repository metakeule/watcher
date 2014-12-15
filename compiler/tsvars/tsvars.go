package tsvars

import (
	"github.com/metakeule/goh4"
	. "github.com/metakeule/watcher/helpers"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type tsVarsCompiler struct {
	*sync.Mutex
	All        []goh4.Class
	Package    string
	Dir        string
	OutputFile string
	Bin        string
	name       string
	Var        string
}

func (ø *tsVarsCompiler) Name() string              { return "tsVars compiler (" + ø.name + ")" }
func (ø *tsVarsCompiler) Dirs() []string            { return []string{ø.Dir} }
func (ø *tsVarsCompiler) Affected(path string) bool { return IsInDir(path, ø.Dir) }
func (ø *tsVarsCompiler) Compile(path string) (output string, err error) {
	// tsvars -in=koelnart/app/backend/class  -target=$GOPATH/src/koelnart/app/backend/typescript/classes.ts
	if ø.Var == "" {
		output, err = Exec(ø.Bin, "-in="+ø.Package, "-target="+ø.OutputFile)
	} else {
		output, err = Exec(ø.Bin, "-in="+ø.Package, "-target="+ø.OutputFile, "-var="+ø.Var)
	}
	return
}

func New(name string, dir string, output string, var_ string) (ø *tsVarsCompiler) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("GOPATH is not set")
	}

	var err error
	ø = &tsVarsCompiler{}
	ø.name = name
	ø.Bin, err = Which("tsvars")
	if err != nil {
		log.Fatal("you do not have tsvars installed. run go get github.com/metakeule/goh4/css/tsvars")
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
