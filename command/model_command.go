package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/volatiletech/inflect"
	"github.com/fengren/gin-scaffold/template"
)

type Field struct {
	FieldType string
	// judge for real tyme is time or something
	Misc string
}

// ModelCommand generates files related to model.
type ModelCommand struct {
	PackageName        string
	ModelName          string
	ModelNamePlural    string
	InstanceName       string
	InstanceNamePlural string
	TemplateName       string
	Fields             map[string]Field
}

// Help prints a help message for this command.
func (command *ModelCommand) Help() {
	fmt.Printf(`Usage:
	gin-scaffold model <model name> <field name>:<field type> ...

Description:
	The gin-scaffold model command creates a new model with the given fields.

Example:
	gin-scaffold model Post Title:string Body:string 
`)
}

func findFieldType(name string) (string, string) {
	misc := ""
	switch name {
	case "text":
		{
			name = "string"
		}
	case "float":
		{
			name = "float64"
		}
	case "boolean":
		{
			name = "bool"
		}
	case "integer":
		{
			name = "int"
		}
	case "decimal":
		{
			name = "int64"
		}
	case
		"time",
		"date",
		"datetime":
		{
			name = "int64"
			misc = "time"
		}
	}

	return name, misc
}

// Converts "<fieldname>:<type>" to {"<fieldname>": "<type>"}
func processFields(args []string) map[string]Field {
	fields := map[string]Field{}
	for _, arg := range args {
		fieldNameAndType := strings.SplitN(arg, ":", 2)
		key := inflect.Titleize(fieldNameAndType[0])
		name, misc := findFieldType(fieldNameAndType[1])
		field:= Field{name, misc}
		fields[key] = field
	}

	return fields
}

// Execute runs this command.
func (command *ModelCommand) Execute(args []string) {
	if len(args) < 2 {
		command.Help()
		os.Exit(2)
	}
	command.ModelName = inflect.Titleize(args[0])
	command.ModelNamePlural = inflect.Pluralize(command.ModelName)

	command.Fields = processFields(args[1:])
	command.InstanceName = inflect.CamelizeDownFirst(command.ModelName)
	command.InstanceNamePlural = inflect.Pluralize(command.InstanceName)
	command.PackageName = template.PackageName()

	outputPath := filepath.Join("models", inflect.Underscore(command.ModelName)+".go")

	// check models folder
	if _, err := os.Stat("models"); os.IsNotExist(err) {
		fmt.Printf("Error: %s\n", err)
		return
	}

	builder := template.NewBuilder("model.go.tmpl")
	builder.WriteToPath(outputPath, command)

	outputPath = filepath.Join("models", inflect.Underscore(command.ModelName)+"_dbsession.go")
	builder = template.NewBuilder("model_dbsession.go.tmpl")
	builder.WriteToPath(outputPath, command)
}
