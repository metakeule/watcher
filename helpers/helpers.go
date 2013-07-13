package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"time"
)

var Debug = false

type execStruct struct {
	output string
	err    error
}

func Exec(name string, opts ...string) (output string, err error) {
	stamp := time.Now().UnixNano()
	id := fmt.Sprintf("%s %s [%v]", name, strings.Join(opts, " "), stamp)
	if Debug {
		log.Printf("executing: %s", id)
	}
	cmd := exec.Command(name, opts...)
	var out []byte
	out, err = cmd.CombinedOutput()
	output = string(out)
	if Debug {
		log.Printf("output of %s:\n%s", id, output)
	}
	if err != nil {
		err = fmt.Errorf(output)
		log.Printf("error of %s:\n%s", id, err.Error())
	}
	return
}

func Which(cmd string) (path string, err error) {
	path, err = Exec("which", cmd)
	path = strings.TrimRight(path, "\n")
	if path == "" {
		err = fmt.Errorf("not found: %s", cmd)
	}
	return
}

func IsIgnored(file string, ignores []string) (is bool, err error) {
	for _, i := range ignores {
		reg, e := regexp.Compile(i)
		if e != nil {
			err = e
			return
		}
		if reg.MatchString(file) {
			is = true
			return
		}
	}
	return
}

func AllAffectedDirs(base string, ignores []string) (all []string, err error) {
	// TODO for every subdir take the affected dirs and merge them
	dirs, err := ioutil.ReadDir(base)
	if err != nil {
		return
	}
	all = []string{}
	for _, dir := range dirs {
		if dir.IsDir() == false {
			continue
		}
		is, e := IsIgnored(dir.Name(), ignores)
		if e != nil {
			err = e
			return
		}
		if is {
			continue
		}
		dirPath := path.Join(base, dir.Name())
		all = append(all, dirPath)
		var subs []string
		subs, err = AllAffectedDirs(dirPath, ignores)
		if err != nil {
			return
		}
		all = append(all, subs...)
	}
	return
}

func CompileToFile(file string, output string) error {
	return ioutil.WriteFile(file, []byte(output), os.FileMode(0644))
}

func IsInDir(path string, dir string) bool {
	return strings.Contains(path, dir)
}
