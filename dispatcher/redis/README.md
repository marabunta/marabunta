# redis backend for marabunta

RPOPLPUSH https://redis.io/commands/rpoplpush#pattern-reliable-queue

* todo
* in progress
* done

# Flow

Each queue is a `LIST`, publishers will use `LPUSH`,  check `BRPOPLPUSH` and `RPOP`

For on demand tasks:

Based on the tasks the scheduler writes to `event_queue`
Agents read from `event_queue` (probably blocking `BRPOPLPUSH`)
Agent gets an event and put in the `processing_queue` when done event is removed from the queue, if it crashes event remains in the `processing_queue` (need to find a GC here)

For scheduled tasks use a `zset`

The dispatcher pools the `zset` and moves the tasks to the `event_queue`

    tasks = ZREVRANGEBYSCORE "new_tasks" <NOW> 0 #this will only take tasks with timestamp lower/equal than now
