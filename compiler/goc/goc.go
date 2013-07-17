package goc

import (
	"fmt"
	. "github.com/metakeule/watcher"
	. "github.com/metakeule/watcher/helpers"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type gocompiler struct {
	*CompileStruct
	GoPath    string
	MainFile  string
	execTests bool
	PidFile   string
}

func NewGoCompiler(mainfile string, pidFile string, dir string, execTests bool, ignore []string, options ...string) Compiler {
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
	return &gocompiler{NewCompiler(dir, "go build", ".go", bin, ignore, options...), gopath, mainfile, execTests, pidFile}
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

	go ø.sighup()

	//go Exec(ø.CompileStruct.Bin, "run", ø.MainFile)
	return
}

func (ø *gocompiler) sighup() {
	b, err := ioutil.ReadFile(ø.PidFile)
	if err != nil {
		// no process, nothing to do
		return
	}
	pid, e := strconv.Atoi(string(b))
	if e != nil {
		panic("can't parse pid " + e.Error())
	}
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -s HUP %v", pid))
	msg, errr := cmd.CombinedOutput()
	if errr != nil {
		log.Printf("could not reload process %v: %s\n", pid, msg)
	}
}
