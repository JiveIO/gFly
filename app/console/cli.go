package main

import (
	_ "app/app/console/commands"  // Autoload commands into pool.
	_ "app/app/console/queues"    // Autoload tasks into queue.
	_ "app/app/console/schedules" // Autoload jobs into schedule.
	"app/core/console"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"os"
)

func main() {
	args := os.Args[1:] // Skip application name

	switch {
	case len(args) > 0 && args[0] == "schedule:run":
		/*---------------------------------------
						Scheduler
		----------------------------------------*/
		// Start scheduler
		console.StartScheduler()
	case len(args) > 0 && args[0] == "queue:run":
		/*---------------------------------------
						QueueJob
		----------------------------------------*/
		// Start queue worker
		console.StartQueueWorker()
	case len(args) > 0 && args[0] == "cmd:run":
		/*---------------------------------------
						Command
		----------------------------------------*/
		// Run command
		console.RunCommands(args[1:])
	}
}
