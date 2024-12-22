package main

import (
	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/pkg/enqueuer"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app_id := os.Getenv("QISCUS_APP_ID")
	secret := os.Getenv("QISCUS_APP_SECRET")
	webHookSecret := os.Getenv("SECRET")
	webhookPort := os.Getenv("PORT")
	amqp := os.Getenv("RABBIT_MQ_DSN")
	maxCustServer, _ := strconv.Atoi(os.Getenv("MAX_SERVED_CUSTOMER_PER_AGENT"))
	agentPoolIntervalSync, _ := strconv.Atoi(os.Getenv("AGENT_POOL_SYNC_INTERVAL"))

	port, _ := strconv.Atoi(webhookPort)

	enq := enqueuer.NewWebHookHandlerEnqueuer(
		port,
		webHookSecret,
		amqp,
		app_id,
		secret,
		maxCustServer,
		agentPoolIntervalSync,
	)
	go enq.RunAllocator()
	enq.RunWebhookHandler()

}
