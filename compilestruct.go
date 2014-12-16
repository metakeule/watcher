package watcher

import (
	. "gopkg.in/metakeule/watcher.v1/helpers"
	"log"
	"path/filepath"
	"sync"
)

type CompileStruct struct {
	*sync.Mutex
	Bin       string
	Options   []string
	Dir       string
	Ext       string
	name      string
	Ignore    []string // ignored directories as regular expressions
	watchDirs []string
}

func NewCompiler(dir string, name string, ext string, bin string, ignore []string, options ...string) (ø *CompileStruct) {
	ø = &CompileStruct{
		Mutex:     &sync.Mutex{},
		Bin:       bin,
		Options:   options,
		Dir:       dir,
		Ext:       ext,
		name:      name,
		Ignore:    ignore,
		watchDirs: []string{},
	}
	var err error
	ø.watchDirs, err = AllAffectedDirs(ø.Dir, ø.Ignore)
	if err != nil {
		log.Fatalf("can't set watching dirs: %s", err.Error())
	}
	ø.watchDirs = append(ø.watchDirs, ø.Dir)
	return
}

func (ø *CompileStruct) Dirs() []string { return ø.watchDirs }

func (ø *CompileStruct) Name() string { return ø.name }

func (ø *CompileStruct) Compile(file string) (output string, err error) {
	opts := append([]string{}, ø.Options...)
	opts = append(opts, file)
	output, err = Exec(ø.Bin, opts...)
	return
}

func (ø *CompileStruct) Affected(file string) bool {
	if filepath.Ext(file) == ø.Ext && IsInDir(file, ø.Dir) {
		return true
	}
	return false
}
