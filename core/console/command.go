package console

import (
	"app/core/log"
	"fmt"
	"regexp"
)

// ===========================================================================================================
// 												Command
// ===========================================================================================================

// ICommand The interface command.
type ICommand interface {
	Validate(parameters CommandParameter) error
	Handle()
}

// Command Abstract command.
type Command struct{}

func (c Command) Validate(parameters CommandParameter) error {
	return nil
}

func (c Command) Handle() {
	log.Panic("Must be implemented from children")
}

// ===========================================================================================================
// 											Command Handler
// ===========================================================================================================

// CommandParameter a map to store parameters
type CommandParameter map[string]interface{}

// Command pool
var commands = make(map[string]ICommand)

// RegisterCommand Register a new command to pool.
func RegisterCommand(name string, cmd ICommand) {
	commands[name] = cmd
}

func RunCommands(args []string) {
	cmdName := args[0]
	cmd, ok := commands[cmdName]

	if !ok {
		fmt.Println("Don't see commands name", cmdName)
		return
	}

	// Parameter pattern `--age=41 --name=John --email="John Land<john@mail.com>" --married=true`
	var parameters = CommandParameter{}
	for _, param := range args[1:] {
		re := regexp.MustCompile(`--(\w+)=([\w*\s+"'<>@./#&-]*)`)
		matches := re.FindStringSubmatch(param)

		parameters[matches[1]] = matches[2]
	}

	err := cmd.Validate(parameters)
	if err != nil {
		fmt.Printf("Error: %v", err)

		return
	}

	cmd.Handle()
}
