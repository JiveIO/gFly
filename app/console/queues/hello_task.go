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
// 					Register task.
// ./artisan queue:run
// ---------------------------------------------------------------

// Auto-register task into queue.
func init() {
	console.RegisterTask(HelloTaskName, HandleHelloTask)
}

// ---------------------------------------------------------------
// 					HelloTask struct.
// ---------------------------------------------------------------

// HelloTaskName Task name.
const HelloTaskName = "hello-world"

// HandleHelloPayload Email delivery payload.
type HandleHelloPayload struct {
	Msg string
}

//----------------------------------------------
// Dispatching functions will make a new Task into Queue.
// A task consists of a type and a payload.
//----------------------------------------------

// DispatchHelloTask Dispatch a new message from somewhere.
func DispatchHelloTask(msg string) (*asynq.Task, error) {
	// Encode Payload
	payload, err := json.Marshal(HandleHelloPayload{Msg: msg})
	if err != nil {
		return nil, err
	}

	log.Infof("Dispatched HelloTask with message %s", msg)

	// Create new task and push to Queue.
	return asynq.NewTask(HelloTaskName, payload), nil
}

// ---------------------------------------------------------------
// Handling functions will pull a new Task from Queue.
// A task consists of a type and a payload.
// Note that it satisfies the asynq.HandlerFunc interface.
//
// Handler doesn't need to be a function. You can define a type
// that satisfies asynq.Handler interface.
// ---------------------------------------------------------------

// HandleHelloTask Handle email delivery task.
func HandleHelloTask(ctx context.Context, t *asynq.Task) error {
	// Decode task Payload
	var p HandleHelloPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Process payload (Just print out)
	log.Infof("Handle HelloTask with message %s", p.Msg)

	return nil
}
