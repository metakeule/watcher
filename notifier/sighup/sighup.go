package sighup

import (
	"fmt"
	"gopkg.in/metakeule/fmtdate.v1"
	"gopkg.in/metakeule/watcher.v1"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type flag int

const (
	_                      = iota
	defaults          flag = 1 << iota
	IgnoreAll              // all errors should be ignored
	IgnoreMissingFile      // missing pid file should be ignored
	IgnoreParseError       // parse error of pid file should be ignored
	IgnoreSighupError      // error while sending sighup signal should be ignored
	IgnoreSuccess          // don't send success message if reload was successfull
)

type sighup struct {
	// the inner notifier
	watcher.Notifier

	// that path of the pid file
	pidFile string

	// ignore flags
	ignores flag

	// start command
	startCommand []string
}

// on error, it only passes the message to the inner notifier
// on success, it parses the pidfile for a pid, calls kill -s HUP
// on the process and reports any error to the inner notifier
// if no error did happen, it passes the original message to inner notifier
// to ignore certain errors, pass the Ignore constants
func New(n watcher.Notifier, pidFile string, startCommand []string, ignoreErrors ...flag) watcher.Notifier {
	ø := &sighup{Notifier: n, pidFile: pidFile, startCommand: startCommand}

	for _, flag := range ignoreErrors {
		ø.ignores = ø.ignores | flag
	}
	return ø
}

// returns if a ignore flag is set
func (ø *sighup) shouldIgnore(f flag) bool {
	if ø.ignores&IgnoreAll != 0 {
		return true
	}
	return ø.ignores&f != 0
}

func (ø *sighup) Error(msg string) {
	ø.Notifier.Error(msg)
}

func (ø *sighup) reportError(format string, i ...interface{}) {
	ø.Notifier.Error(fmt.Sprintf(format, i...))
}

func (ø *sighup) readUntilPidChanged(retries int, pid int) (int, error) {
	if retries > 10 {
		return 0, fmt.Errorf("no process running")
	}
	time.Sleep(time.Millisecond * 500)
	b, err := ioutil.ReadFile(ø.pidFile)
	if err != nil {
		return ø.readUntilPidChanged(retries+1, pid)
	}
	newpid, e := strconv.Atoi(string(b))
	if e != nil {
		return 0, e
	}
	if newpid == pid {
		return ø.readUntilPidChanged(retries+1, pid)
	}
	return newpid, nil
}

func (ø *sighup) hasStartCommand() bool {
	return len(ø.startCommand) > 0
}

var (
	hexRegex = regexp.MustCompile(`\+{0,1}0x[0-9a-h]+`)
	gofile   = regexp.MustCompile(`(\.go:\-{0,1}[0-9]+)\s*$`)
)

func filterLogs(in string) (out string) {
	str := fmtdate.Format("YYYY/MM/DD", time.Now())
	quoted := regexp.QuoteMeta(str)
	reg := regexp.MustCompile("^" + quoted)
	var o []string

	for _, line := range strings.Split(in, "\n") {
		if !reg.MatchString(line) {
			hexreplaced := hexRegex.ReplaceAllString(line, "")
			repl := gofile.ReplaceAllString(hexreplaced, "$1 «\n")
			o = append(o, repl)
		}
	}

	return strings.Join(o, "\n")
}

func (ø *sighup) runStart() {
	if ø.hasStartCommand() {
		// fmt.Printf("try to start: %#v\n", strings.Join(ø.startCommand, " "))
		cmd := exec.Command(ø.startCommand[0], ø.startCommand[1:]...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			// fmt.Printf("Can't start process: %#v\n", string(out))
			ø.reportError(filterLogs(string(out)))
			return
		}
		// fmt.Println("process started")

	}
}

// any errors are reported to the inner notifier
func (ø *sighup) Success(msg string) {
	ø.Notifier.Success(msg)
	b, err := ioutil.ReadFile(ø.pidFile)
	if err != nil {
		if ø.hasStartCommand() {
			if !ø.shouldIgnore(IgnoreSuccess) {
				ø.Notifier.Success("started process")
			}
			go ø.runStart()
			return
		}
		if !ø.shouldIgnore(IgnoreMissingFile) {
			ø.reportError("Can't sighup, PID file %#v missing", ø.pidFile)
		}
		return
	}
	pid, e := strconv.Atoi(string(b))
	if e != nil {
		if !ø.shouldIgnore(IgnoreParseError) {
			ø.reportError("Can't parse pid to int: %#v", e.Error())
		}
		return
	}
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -s HUP %v", pid))
	out, errr := cmd.CombinedOutput()
	if errr != nil {
		if !ø.shouldIgnore(IgnoreSighupError) {
			ø.reportError("could not reload process %v: %s\n", pid, out)
		}
		return
	}

	newpid, pidErr := ø.readUntilPidChanged(1, pid)
	if pidErr != nil {
		if !ø.shouldIgnore(IgnoreSighupError) {
			ø.reportError("could not reload process %v: %s\n", pid, pidErr.Error())
		}
		return
	}

	_, procErr := os.FindProcess(newpid)
	if procErr != nil {
		if !ø.shouldIgnore(IgnoreSighupError) {
			ø.reportError("could not reload process %v => %v (process not restarted: %s\n", pid, newpid, procErr.Error())
		}
		return
	}

	if !ø.shouldIgnore(IgnoreSuccess) {
		ø.Notifier.Success("reloaded")
	}
}
