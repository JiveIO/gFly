package console

import (
	"app/core/errors"
	"app/core/log"
	"app/core/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"sync"
	"time"
)

// ===========================================================================================================
// 												Queue task
// ===========================================================================================================

// ITask The interface task.
type ITask interface {
	// Dequeue get out and process task in queue.
	Dequeue(ctx context.Context, t *asynq.Task) error
}

// Task Abstract task.
type Task struct{}

func (t Task) Dequeue(ctx context.Context, task *asynq.Task) error {
	return errors.NotYetImplemented
}

// ===========================================================================================================
// 											Queue handler
// ===========================================================================================================

func getRedisClientOpt() asynq.RedisClientOpt {
	// Build Redis connection URL.
	redisConnURL := fmt.Sprintf(
		"%s:%d",
		utils.Getenv("REDIS_HOST", "localhost"),
		utils.Getenv("REDIS_PORT", 6379),
	)

	// Define Redis database number.
	dbNumber := utils.Getenv("REDIS_DB_NUMBER", 0)

	return asynq.RedisClientOpt{
		Addr:     redisConnURL,
		Password: utils.Getenv("REDIS_PASSWORD", ""),
		DB:       dbNumber,
	}
}

var client = asynq.NewClient(getRedisClientOpt())

// StartQueueWorker Start queue worker.
// Worker handles a Task(job) was pushed to Queue (Redis) from somewhere.
//   - Many queue Workers run (in a separated server) and try to get Task in Queue to handle
//   - Somewhere in or outside application add a new Task to Queue.
func StartQueueWorker() {
	// Create queue worker
	srv := asynq.NewServer(
		getRedisClientOpt(),
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optional specify multiple queues (Queue type) with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()

	// Register task handlers
	for key, task := range queueTasks {
		// Register a task handler
		mux.HandleFunc(key, task.Dequeue)
		log.Infof("Init queue task %s", key)
	}

	// Start queue worker
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

// queueTasks pool to store task in queue.
var queueTasks = make(map[string]ITask)

// RegisterTask Register a new task to pool.
func RegisterTask(task ITask, name string) {
	queueTasks[name] = task
}

// DispatchTask push a task to queue.
func DispatchTask(data interface{}, name string) {
	startTime := time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := handleEnqueue(data, name)
		if err != nil {
			log.Errorf("Error %v", err)
		}
	}()

	wg.Wait()

	log.Infof("[RUN] Dispatch Task %s - %v", name, time.Since(startTime))
}

func handleEnqueue(data interface{}, name string) error {
	// Get the corresponding task instance to process
	_, ok := queueTasks[name]
	if !ok {
		log.Errorf("Invalid queue `%s`", name)

		return errors.InvalidParameter
	}

	// Encode Payload
	payload, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Encode error %v. Error %v", data, err)

		return err
	}

	// Create new task and push to Queue.
	task := asynq.NewTask(name, payload)

	// Enqueue task
	info, err := client.Enqueue(task)
	if err != nil {
		log.Errorf("Could not enqueue task: %v", err)

		return err
	}
	log.Infof("Enqueue task: %v", info)

	return nil
}
