package goc

import (
	. "github.com/metakeule/watcher"
	. "github.com/metakeule/watcher/helpers"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type gocompiler struct {
	*CompileStruct
	GoPath    string
	MainFile  string
	execTests bool
}

func NewGoCompiler(mainfile string, dir string, execTests bool, ignore []string, options ...string) Compiler {
	if mainfile != "" {
		if filepath.Ext(mainfile) != ".go" {
			log.Fatalf("mainfile %#v is not a go file", mainfile)
		}
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("GOPATH is not set")
	}

	bin, err := Which("go")
	if err != nil {
		log.Fatal("you do not have the go compiler installed.")
	}
	/*
	   if len(options) == 0 {
	       // recommended options
	       //options := []string{"--strict-imports", "--verbose", "--no-color", "--line-numbers=all", "--strict-math=on", "--strict-units=off"}
	   }
	*/
	return &gocompiler{NewCompiler(dir, "go build", ".go", bin, ignore, options...), gopath, mainfile, execTests}
}

func (ø *gocompiler) Compile(file string) (output string, err error) {
	if file != ø.MainFile {
		dir := filepath.Dir(file)
		packageName := strings.Replace(dir, path.Join(ø.GoPath, "src")+"/", "", 1)
		output, err = Exec(ø.CompileStruct.Bin, "install", packageName)
		if err != nil {
			return
		}

		if ø.execTests {
			output, err = Exec(ø.CompileStruct.Bin, "test", packageName)
			if err != nil {
				return
			}
		}
	}

	if ø.MainFile == "" {
		return
	}

	output, err = Exec(ø.CompileStruct.Bin, "build", ø.MainFile)
	if err != nil {
		return
	}

	go Exec(ø.CompileStruct.Bin, "run", ø.MainFile)
	return
}
