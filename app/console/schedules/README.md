# Scheduler `app/console/schedules/`

Scheduler only handle a instance which implemented interface `IJob`

```go
// IJob The interface task.
type IJob interface {
    GetTime() string // Schedule time. Example "0/5 * * * * *"
    Handle() // Method to handle the job
}
```

Let example you want to create a job name `TickTickJob`. A struct `TickTickJob` will be created 
```go
// TickTickJob struct for ticktick job.
type TickTickJob struct{}
```

You want the job will be called every 5 seconds. So, need to implement method `GetTime() string` of interface `IJob` for struct `TickTickJob`

```go
// GetTime Get time format.
func (c *TickTickJob) GetTime() string {
    return "0/5 * * * * *"
}
```
The handler of job will make a random message from pool (have 3 messages) and display.


```go
// Handle Process the job.
func (c *TickTickJob) Handle() {
    var messages = []string{
        "Today is a good day to learn gFly framework",
        "What is the weather today?... Umm",
        "The job will be scheduled. The task will be enqueued!",
    }

    // Get random Index
    idx := utils.RandInt64(int64(len(messages)))

    log.Infof("Message %s", messages[idx])
    log.Infof("TickTickJob :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

The latest important thing is register job to scheduler

```go
// Auto-register job into scheduler.
func init() {
    console.RegisterJob(&TickTickJob{})
}
```

## Job file `ticktick_job.go`

```go
package schedules

import (
    "app/core/console"
    "app/core/log"
    "app/core/utils"
    "time"
)

// ---------------------------------------------------------------
//                      Register job.
// ---------------------------------------------------------------

// Auto-register job into scheduler.
func init() {
    console.RegisterJob(&TickTickJob{})
}

// ---------------------------------------------------------------
//                      TickTickJob struct.
// ---------------------------------------------------------------

// TickTickJob struct for ticktick job.
type TickTickJob struct{}

// GetTime Get time format.
func (c *TickTickJob) GetTime() string {
    return "0/5 * * * * *"
}

// Handle Process the job.
func (c *TickTickJob) Handle() {
    var messages = []string{
        "Today is a good day to learn gFly framework",
        "What is the weather today?... Umm",
        "The job will be scheduled. The task will be enqueued!",
    }

    // Get random Index
    idx := utils.RandInt64(int64(len(messages)))

    log.Infof("Message %s", messages[idx])
    log.Infof("TickTickJob :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
```

**Note:** Put the job file `ticktick_job.go` correct folder `app/console/schedules/`

## Check `TickTickJob` job

Run scheduler

```bash
./build/artisan schedule:run
```