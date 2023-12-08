package console

import (
	"app/core/log"
	"app/core/utils"
	"context"
	"fmt"
	"github.com/hibiken/asynq"
)

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

// TODO Need to close `defer client.Close()`.
var client = asynq.NewClient(getRedisClientOpt())

// QueueClient Create queue client
func QueueClient() *asynq.Client {
	return client
}

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
	for key, handler := range queueTasks {
		// Register a task handler
		mux.HandleFunc(key, handler)
		log.Infof("Init queue task %s", key)
	}

	// Start queue worker
	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

type TaskHandler func(ctx context.Context, t *asynq.Task) error

// queueTasks pool to store task in queue.
var queueTasks = make(map[string]TaskHandler)

// RegisterTask Register a new task to pool.
func RegisterTask(taskName string, task TaskHandler) {
	queueTasks[taskName] = task
}
