package commands

import (
	"app/app/console/queues"
	"app/core/console"
	"app/core/log"
	"time"
)

// ---------------------------------------------------------------
// 					Register command.
// ./artisan cmd:run hello-world
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(HelloCommandName, &HelloCommand{})
}

// ---------------------------------------------------------------
// 					HelloCommand struct.
// ---------------------------------------------------------------

// HelloCommandName Define command name
const HelloCommandName = "hello-world"

// HelloCommand struct for hello command.
type HelloCommand struct {
	console.Command
}

// Handle Process command.
func (c *HelloCommand) Handle() {
	// Dispatch a task into Queue.
	_, err := queues.DispatchHelloTask("Hello world")
	if err != nil {
		log.Errorf("Can not dispatch HelloTask %v", err)
	}

	log.Infof("HellCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
