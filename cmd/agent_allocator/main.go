package main

import (
	"fmt"
	"log"
	"os"
	"time"

	agentpool "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/pkg/agent_pool"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app_id := os.Getenv("QISCUS_APP_ID")
	secret := os.Getenv("QISCUS_APP_SECRET")

	agentPoolCfg := agentpool.AgentPoolConfig{
		MaxServedCustomerPerAgent: 2,
		SyncInterval:              300,
		QiscusApiAuthAppId:        app_id,
		QiscusApiAuthSecret:       secret,
		QiscusBaseHttpApiHost:     "https://multichannel.qiscus.com",
	}

	agentPool := agentpool.NewAgentPoolAllocator(agentPoolCfg)
	fmt.Println(agentPool)
	fmt.Println("MAIN ALLOCATE 1", agentPool.AllocateAgent(285685151))
	fmt.Println("MAIN ALLOCATE 2", agentPool.AllocateAgent(285685197))
	fmt.Println("MAIN ALLOCATE 3", agentPool.AllocateAgent(285685241))

	time.Sleep(5 * time.Hour)

}
