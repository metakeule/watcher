package app2

import (
	"github.com/metakeule/watcher"
	"github.com/metakeule/watcher/compiler/cssvars"
	"github.com/metakeule/watcher/compiler/less"
	"github.com/metakeule/watcher/compiler/tsvars"
	"github.com/metakeule/watcher/compiler/typescript"
	"path"
)

func Compilers(baseDir string, app string) []watcher.Compiler {
	appDir := path.Join(baseDir, "app", app)
	staticDir := path.Join(baseDir, "static")

	lessDir := path.Join(appDir, "compiler", "css")
	lessOutput := path.Join(staticDir, "css", app, "all.css")
	lessMain := path.Join(lessDir, "main.less")

	typeScriptDir := path.Join(appDir, "compiler", "js")
	typeScriptOutput := path.Join(staticDir, "js", app)

	classDir := path.Join(appDir, "compiler", "class")
	classLessOutput := path.Join(lessDir, "class.less")
	classTsOutput := path.Join(typeScriptDir, "class.ts")

	idDir := path.Join(appDir, "compiler", "id")
	idLessOutput := path.Join(lessDir, "id.less")
	idTsOutput := path.Join(typeScriptDir, "id.ts")

	return []watcher.Compiler{
		cssvars.New(app+" class.less", classDir, classLessOutput, ""),
		cssvars.New(app+" id.less", idDir, idLessOutput, ""),
		tsvars.New(app+" class.ts", classDir, classTsOutput, ""),
		tsvars.New(app+" id.ts", idDir, idTsOutput, ""),
		less.New(lessMain, lessDir, lessOutput, []string{}),
		typescript.New(typeScriptDir, typeScriptOutput, []string{}),
	}
}
