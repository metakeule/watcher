package less

import (
	"github.com/metakeule/watcher"
	. "github.com/metakeule/watcher/helpers"
	"log"
)

type less struct {
	*watcher.CompileStruct
	OutPutFile string
	MainFile   string
}

func New(mainFile string, dir string, outputFile string, ignore []string, options ...string) watcher.Compiler {
	bin, err := Which("lessc")
	if err != nil {
		log.Fatal("you do not have the less compiler installed. please run 'npm install -g less'")
	}
	if len(options) == 0 {
		// recommended options 
		options = []string{"--strict-imports", "--verbose", "--no-color", "--line-numbers=all", "--strict-math=on", "--strict-units=off"}
	}
	return &less{watcher.NewCompiler(dir, "less (lessc)", ".less", bin, ignore, options...), outputFile, mainFile}
}

func (ø *less) Compile(file string) (output string, err error) {
	opts := append([]string{}, ø.CompileStruct.Options...)
	opts = append(opts, ø.MainFile)
	output, err = Exec(ø.Bin, opts...)
	if err != nil {
		return
	}
	err = CompileToFile(ø.OutPutFile, output)
	return
}
