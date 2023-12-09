package console

import (
	"app/core/log"
	"app/core/utils"
	"github.com/robfig/cron/v3"
)

// IJob The interface task.
type IJob interface {
	GetTime() string
	Handle()
}

// Job pool
var jobs []IJob

// RegisterJob Register a new job to pool.
func RegisterJob(job IJob) {
	jobs = append(jobs, job)
}

func StartScheduler() {
	c := cron.New(cron.WithSeconds())

	// Define the Cron job schedule
	for _, job := range jobs {
		_, err := c.AddFunc(job.GetTime(), job.Handle)
		if err != nil {
			return
		}
		log.Infof("Init schedule job %s", utils.ReflectType(job))
	}

	// Start the Cron job scheduler
	c.Start()

	// Run forever
	var forever chan struct{}
	<-forever
}
