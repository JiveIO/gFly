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
	console.RegisterCommand(&HelloCommand{}, "hello-world")
}

// ---------------------------------------------------------------
// 					HelloCommand struct.
// ---------------------------------------------------------------

// HelloCommand struct for hello command.
type HelloCommand struct {
	console.Command
}

// Handle Process command.
func (c *HelloCommand) Handle() {
	// Dispatch a task into Queue.
	console.DispatchTask(queues.NewHelloTask("Hello"))

	log.Infof("HellCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
