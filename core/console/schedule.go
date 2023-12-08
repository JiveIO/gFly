package console

import (
	"app/core/log"
	"github.com/robfig/cron/v3"
)

// IJob The interface task.
type IJob interface {
	GetTime() string
	Handle()
}

// Job pool
var jobs = make(map[string]IJob)

// RegisterJob Register a new job to pool.
func RegisterJob(name string, job IJob) {
	jobs[name] = job
}

func StartScheduler() {
	c := cron.New(cron.WithSeconds())

	// Define the Cron job schedule
	for name, job := range jobs {
		_, err := c.AddFunc(job.GetTime(), job.Handle)
		if err != nil {
			return
		}
		log.Infof("Init schedule job %s", name)
	}

	// Start the Cron job scheduler
	c.Start()

	// Run forever
	var forever chan struct{}
	<-forever
}
