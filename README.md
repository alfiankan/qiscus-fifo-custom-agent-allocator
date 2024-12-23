# QISCUS Custom Agent Allocator RoundRobin


## Agent Allocation Strategy
Allocator will allocate to agent that have `n` max served per agent, and messages (room) are served as FIFO (First in First out) using roundrobin strategy checking agent to agent.

## System Design



## Dependencies
- Golang
- RabbitMQ
- Qiscus API


## How to run
- ensure you have rabbitmq instance to be accessed
- setup `.env` file
    ```bash
        QISCUS_APP_ID=xxxx
        QISCUS_APP_SECRET=xxxx
        QISCUS_API_HOST=https://multichannel.qiscus.com
        PORT=5005
        SECRET=12345
        RABBIT_MQ_DSN=amqp://qiscus:qiscus@localhost:5672
        MAX_SERVED_CUSTOMER_PER_AGENT=2
        AGENT_POOL_SYNC_INTERVAL=120
        QUEUE_BACKOFF_SLEEP_INTERVAL_SECOND=120
    ```

- run `go run ./cmd/agent_allocator/...`
- add `<host>/allocate?secret=12345` to custom agent allocator webhook setting
- Test via qiscus multichannel.
