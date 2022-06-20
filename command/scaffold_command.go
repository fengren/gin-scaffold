package command

import (
	"fmt"
	"os"

        "github.com/volatiletech/inflect"
)

type ScaffoldCommand struct {
}

func (command *ScaffoldCommand) Help() {
	fmt.Printf(`Usage:
	gin-scaffold scaffold <controller name> <field name>:<field type> ...

Description:
	The gin-scaffold scaffold command creates a new controller and model with the given fields.

Example:
	gin-scaffold controller Post Title:string Body:string
`)
}

func (command *ScaffoldCommand) Execute(args []string) {
	if len(args) == 0 {
		command.Help()
		os.Exit(2)
	}
	args[0] = inflect.Singularize(args[0])
	modelCommand := &ModelCommand{}
	modelCommand.Execute(args)

	controllerCommand := &ControllerCommand{}
	controllerCommand.Execute([]string{modelCommand.ModelNamePlural})
}
