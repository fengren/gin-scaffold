package command

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fengren/gin-scaffold/template"
)

var (
	dirsToCreate = []string{
		"models",
		"controllers",
		"helpers",
		"config",
		"templates",
	}
)

type InitCommand struct {
	ProjectDir         string
	ProjectName        string
	DatabaseNamePrefix string
	PackageName        string
}

func (command *InitCommand) Help() {
	fmt.Printf(`Usage:
	gin-scaffold init <app path>

Description:
	The gin-scaffold init command creates a new gin application.

Example:
	gin-scaffold init blog
`)
}

func (command *InitCommand) Execute(args []string) {
	if len(args) == 0 {
		command.Help()
		os.Exit(2)
	}
	projectDir, err := filepath.Abs(args[0])
	if err != nil {
		panic(err)
	}

	wd, _ := os.Getwd()
	wd = filepath.ToSlash(wd)
	root := ""
	for _, p := range filepath.SplitList(os.Getenv("GOPATH")) {
		p = filepath.ToSlash(p)
		if strings.HasPrefix(strings.ToLower(wd), strings.ToLower(filepath.ToSlash(filepath.Join(p, "src"))+"/")) {
			root = wd[len(p+"/src/"):]
		}
	}

	command.ProjectName = filepath.Base(projectDir)
	command.ProjectDir = projectDir
	command.DatabaseNamePrefix = filepath.Base(projectDir)
	command.PackageName = path.Join(root, command.ProjectName)
	command.createLayout()

	command.installFiles("helpers")
	command.installFiles("config")
	command.installFiles("controllers")

	command.installFile("", "main.go.tmpl", command.ProjectName+".go")
}

func (command *InitCommand) installFiles(dirName string) {
	helperFiles, err := filepath.Glob(template.TemplatePath(filepath.Join(dirName, "*.tmpl")))
	if err != nil {
		panic(err)
	}

	for _, templateFile := range helperFiles {
		outputFileName := filepath.Base(templateFile)
		outputFileName = strings.TrimRight(outputFileName, ".tmpl")
		command.installFile(dirName, templateFile, outputFileName)
	}
}

func (command *InitCommand) installFile(dirName string, templateFile string, outputFileName string) {
	builder := template.NewBuilder(templateFile)
	builder.WriteToPath(filepath.Join(command.ProjectDir, dirName, outputFileName), command)
}

func (command *InitCommand) directoryInRoot(path string) string {
	return filepath.Join(command.ProjectDir, path)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (command *InitCommand) createLayout() {
	for _, dirName := range dirsToCreate {
		path := command.directoryInRoot(dirName)
		must(os.MkdirAll(path, 00755))
	}
}
