package typescript

import (
	"github.com/metakeule/watcher"
	. "github.com/metakeule/watcher/helpers"
	"log"
)

type typescript struct {
	*watcher.CompileStruct
	OutPutFile string
}

/*
func (ø *typescript) Compile(file) (output string, err error) {
    output, err = ø.compiler.Compile(file)
    if err != nil {
        return
    }
    err = CompileToFile(ø.OutPutFile, output)
    return
}
*/

func New(dir string, outputDir string, ignore []string, options ...string) watcher.Compiler {
	bin, err := Which("tsc")
	if err != nil {
		log.Fatal("you do not have the typescript compiler installed. please run 'npm install -g typescript'")
	}
	if len(options) == 0 {
		// recommended options 
		options = []string{"--out", outputDir, "-c", "--disallowbool", "--disallowbool", "--sourcemap", "--module", "amd"}
	}
	return &typescript{watcher.NewCompiler(dir, "typescript (tsc)", ".ts", bin, ignore, options...), outputDir}
}
