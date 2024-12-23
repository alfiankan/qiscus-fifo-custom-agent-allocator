package main

import (
	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/pkg/enqueuer"
	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
)

func main() {

	cfg := utils.LoadApplicationConfig()

	enq := enqueuer.NewWebHookHandlerEnqueuer(cfg)
	go enq.RunAllocator()
	enq.RunWebhookHandler()
}
