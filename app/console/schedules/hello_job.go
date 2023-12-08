package schedules

import (
	"app/core/console"
	"app/core/log"
	"time"
)

// ---------------------------------------------------------------
// 					Register job.
// ---------------------------------------------------------------

// Auto-register task into queue.
func init() {
	console.RegisterJob(HelloJobName, &HelloJob{})
}

// ---------------------------------------------------------------
// 					HelloJob struct.
// ---------------------------------------------------------------

// HelloJobName Define job name
const HelloJobName = "hello-world"

// HelloJob struct for hello job.
type HelloJob struct{}

// GetTime Get time format.
func (c *HelloJob) GetTime() string {
	return "0/2 * * * * *"
}

// Handle Process the job.
func (c *HelloJob) Handle() {
	log.Infof("HelloJob :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
