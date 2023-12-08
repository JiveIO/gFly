# Folder app/console/queues/

Worker handles a Task was pushed to Queue (Redis) from somewhere.
- Many queue Workers run (in a separated server) and try to get Task in Queue to handle.
- Somewhere in or outside application add a new Task to Queue.
