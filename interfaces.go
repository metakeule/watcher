package watcher

type (
	Compiler interface {
		Compile(path string) (string, error)
		Name() string
		Lock()
		Unlock()
		// directories to watch
		Dirs() []string
		Affected(path string) bool
	}

	Notifier interface {
		Error(msg string)
		Success(msg string)
	}
)
