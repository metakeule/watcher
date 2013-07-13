package app

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

	lessDir := path.Join(appDir, "less")
	lessOutput := path.Join(staticDir, "css", app, "all.css")
	lessMain := path.Join(lessDir, "main.less")

	typeScriptDir := path.Join(appDir, "typescript")
	typeScriptOutput := path.Join(staticDir, "js", app)

	classDir := path.Join(appDir, "class")
	classLessOutput := path.Join(lessDir, "classes.less")
	classTsOutput := path.Join(typeScriptDir, "Classes.ts")

	idDir := path.Join(appDir, "id")
	idLessOutput := path.Join(lessDir, "ids.less")
	idTsOutput := path.Join(typeScriptDir, "Ids.ts")

	return []watcher.Compiler{
		cssvars.New(app+" classes.less", classDir, classLessOutput, ""),
		cssvars.New(app+" ids.less", idDir, idLessOutput, ""),
		tsvars.New(app+" Classes.ts", classDir, classTsOutput, ""),
		tsvars.New(app+" Ids.ts", idDir, idTsOutput, ""),
		less.New(lessMain, lessDir, lessOutput, []string{}),
		typescript.New(typeScriptDir, typeScriptOutput, []string{}),
	}
}
