package watcher

import (
	"code.google.com/p/go.exp/fsnotify"
	"log"
	"os"
	"sync"
	"time"
)

type ProjectWatcher struct {
	*sync.Mutex
	Compilers []Compiler
	Watcher   *fsnotify.Watcher
	Notifier  Notifier
	Pool      map[Compiler]string
	Frequency time.Duration
	Ready     chan int
}

func New(notifier Notifier, compilers ...Compiler) (ø *ProjectWatcher) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("can't create watcher: %s\n", err.Error())
	}
	return &ProjectWatcher{
		Mutex:     &sync.Mutex{},
		Compilers: compilers,
		Watcher:   watcher,
		Ready:     make(chan int, 1),
		Notifier:  notifier,
		Frequency: time.Millisecond * 20,
		// each compiler should be handled at one time with the last file succeeding
		Pool: map[Compiler]string{},
	}
}

func (ø *ProjectWatcher) SendMessages() {
	for {
		ø.Lock()
		for comp, file := range ø.Pool {
			log.Printf("handle %s with %s\n", file, comp.Name())
			out, err := comp.Compile(file)
			if err != nil {
				ø.Notifier.Error(out)
			} else {
				ø.Notifier.Success("compiled " + file)
			}
			delete(ø.Pool, comp)
		}
		ø.Unlock()
		time.Sleep(ø.Frequency)
	}
}

func (ø *ProjectWatcher) HandleFile(path string) {
	// log.Printf("trying to handle: %s\n", path)
	for _, c := range ø.Compilers {
		c.Lock()
		if c.Affected(path) {
			ø.Pool[c] = path

		}
		c.Unlock()
	}
}

func (ø *ProjectWatcher) Run() (err error) {
	for _, c := range ø.Compilers {
		for _, d := range c.Dirs() {
			err = ø.Watcher.Watch(d)
			if err != nil {
				return
			}
		}
	}

	go ø.SendMessages()

	go func() {
		for {
			select {
			case ev := <-ø.Watcher.Event:
				//log.Println("event: (create:%v)", ev, ev.IsCreate())

				what := ""
				handleIt := true
				switch {
				case ev.IsCreate():
					what = "created"
					d, err := os.Stat(ev.Name)
					if err == nil {
						if d.IsDir() {
							ø.Lock()
							log.Println("added ", ev.Name, "- start watching")
							ø.Watcher.Watch(ev.Name)
							ø.Unlock()
						}
					}
					handleIt = false
				case ev.IsDelete():
					handleIt = false
					what = "deleted"
					ø.Lock()
					ø.Watcher.RemoveWatch(ev.Name)
					ø.Unlock()

				case ev.IsModify():
					what = "modified"
				case ev.IsRename():
					handleIt = false
					what = "renamed"
				}
				_ = what
				//log.Println("file: ", ev.Name, " ", what)

				if handleIt {
					ø.Lock()
					ø.HandleFile(ev.Name)
					ø.Unlock()
				}

			case err := <-ø.Watcher.Error:
				log.Println("watcher error:", err)
			}
		}
		ø.Lock()
		ø.Ready <- 1
		ø.Unlock()
	}()
	return
}

/*
   Flags
       FSN_CREATE = 1
       FSN_MODIFY = 2
       FSN_DELETE = 4
       FSN_RENAME = 8

       FSN_ALL = FSN_MODIFY | FSN_DELETE | FSN_RENAME | FSN_CREATE
*/
//watcher.WatchFlags(ev.Name string, flags uint32) error {

// watcher.RemoveWatch(path string)

//watcher.Close()
