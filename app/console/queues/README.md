# Queue `app/console/queues/`

Queue handler only handle a instance which implemented interface `ITask`

```go
// ITask The interface task.
type ITask interface {
    // Dequeue get out and process task in queue.
    Dequeue(ctx context.Context, t *asynq.Task) error
}
```

Let example you want to create a task name `PingTask`. Not required but you should embed `console.Task` inside `PingTask` task. A struct `PingTask` will be created

```go
// PingTask Ping task.
type PingTask struct {
    console.Task
}
```

The task will have a task name `ping` and `PingTaskPayload` struct as the payload will be processed by task.

```go
// PingTaskPayload Task payload.
type PingTaskPayload struct {
    Message string
}
```

Should make a PingTask constructor function to be used somewhere. The function help somewhere want to dispatch the task `PingTask` easier.

```go
// NewPingTask Constructor PingTask.
func NewPingTask(message string) (PingTaskPayload, string) {
    return PingTaskPayload{
        Message: message,
    }, "ping"
}
```

So, when you want to push a new task to queue. Just make add below code

```go
console.DispatchTask(queues.NewPingTask("Hello"))
```

So, you dispatched a task to queue. Now need a method to handle the `task` in queue.

```go
// Dequeue Handle a task in queue.
func (t PingTask) Dequeue(ctx context.Context, task *asynq.Task) error {
    // Decode task payload
    var payload PingTaskPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }

    // Process payload
    log.Infof("Handle PingTask with message %s", payload.Message)

    return nil
}
```

The latest important thing is register task

```go
// Auto-register task into queue.
func init() {
    console.RegisterTask(&PingTask{}, "ping")
}
```

## Job file `ping_task.go`

```go
package queues

import (
    "app/core/console"
    "app/core/log"
    "context"
    "encoding/json"
    "fmt"
    "github.com/hibiken/asynq"
)

// ---------------------------------------------------------------
//                          Register task.
// ---------------------------------------------------------------

// Auto-register task into queue.
func init() {
    console.RegisterTask(&PingTask{}, "ping")
}

// ---------------------------------------------------------------
//                          Task info.
// ---------------------------------------------------------------

// NewPingTask Constructor PingTask.
func NewPingTask(message string) (PingTaskPayload, string) {
    return PingTaskPayload{
        Message: message,
    }, "ping"
}

// PingTaskPayload Task payload.
type PingTaskPayload struct {
    Message string
}

// PingTask Hello task.
type PingTask struct {
    console.Task
}

// Dequeue Handle a task in queue.
func (t PingTask) Dequeue(ctx context.Context, task *asynq.Task) error {
    // Decode task payload
    var payload PingTaskPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
    }

    // Process payload
    log.Infof("Handle PingTask with message `%s`", payload.Message)

    return nil
}
```

**Note:** Put the task `ping_task.go` correct folder `app/console/queues/`

## Check `PingTask` task

Somewhere in your code just add below code to push a new task to queue.

```bash
console.DispatchTask(queues.NewPingTask("Ping ..."))
```

**TIP:** You could add this code to `app/console/commands/hello_command.go` and run `./build/artisan cmd:run hello-world` for testing ;).

Run command (from another terminal)
```bash
./build/artisan queue:run
```