package agentpool

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
)

type Agent struct {
	ID            int
	IsAvailable   bool
	CustomerCount int
}

type AgentPoolConfig struct {
	MaxServedCustomerPerAgent int
	SyncInterval              int
	QiscusBaseHttpApiHost     string
	QiscusApiAuthAppId        string
	QiscusApiAuthSecret       string
}

type AgentPool struct {
	agents                []Agent
	config                AgentPoolConfig
	logLabel              string
	lock                  sync.RWMutex
	lastRoundRobinPointer int
	qiscusApi             QiscusApiClient
}

func NewAgentPoolAllocator(cfg AgentPoolConfig) *AgentPool {
	pool := &AgentPool{
		config:                cfg,
		agents:                []Agent{},
		logLabel:              "AGENT_ALLOCATOR",
		lastRoundRobinPointer: 0,
		qiscusApi:             *NewQiscusApiClient(cfg.QiscusBaseHttpApiHost, cfg.QiscusApiAuthAppId, cfg.QiscusApiAuthSecret),
	}
	pool.syncAgent()
	go pool.startBackgroundSync()

	return pool
}

func (self *AgentPool) startBackgroundSync() {
	utils.LogWrite(self.logLabel, utils.LOG_INFO, "Starting Background Sync Ticker")
	ticker := time.NewTicker(time.Duration(self.config.SyncInterval) * time.Second)

	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			utils.LogWrite(self.logLabel, utils.LOG_INFO, "BACKGROUND SYNC TICKED", t.String())
			self.syncAgent()
		}
	}
}

// remove anavailable agent and get new agent if exist on api
func (self *AgentPool) syncAgent() {

	syncedAgent := []Agent{}

	pageNow := 1

	for {
		agentsResponseData, err := self.qiscusApi.GetAllAgents(pageNow, 10)
		if err != nil {
			break
		}
		if len(agentsResponseData.Data.Agents.Data) == 0 {
			break
		}

		for _, agent := range agentsResponseData.Data.Agents.Data {
			syncedAgent = append(syncedAgent, Agent{
				ID:            agent.ID,
				IsAvailable:   false,
				CustomerCount: agent.CurrentCustomerCount,
			})
		}
		pageNow += 1
	}

	// acquire lock to update pool
	self.lock.Lock()
	self.agents = syncedAgent
	self.lock.Unlock()
	utils.LogWrite(self.logLabel, utils.LOG_INFO, "AGENT ARE SYNCED", fmt.Sprintf("Total Agent: %d", len(self.agents)))

	return
}

func (self *AgentPool) AllocateAgent(roomId int) (err error) {
	utils.LogWrite(self.logLabel, utils.LOG_INFO, fmt.Sprintf("Allocating to room: %d", roomId))

	self.lock.Lock()

	// reset round robin if more than offset
	if len(self.agents) > self.lastRoundRobinPointer {
		self.lastRoundRobinPointer = 0
	}

	for {
		if self.lastRoundRobinPointer > len(self.agents)-1 {
			break
		}
		pickedAgent := self.agents[self.lastRoundRobinPointer]

		// get latest customer served count from qiscusApi or continue to next agent ids in pool
		qiscusAgent, err := self.qiscusApi.GetAgentDetailById(pickedAgent.ID)
		if err != nil {
			self.lastRoundRobinPointer += 1
			continue
		}
		pickedAgent.CustomerCount = qiscusAgent.CurrentCustomerCount
		pickedAgent.IsAvailable = qiscusAgent.IsAvailable

		if pickedAgent.CustomerCount == self.config.MaxServedCustomerPerAgent || !pickedAgent.IsAvailable {
			self.lastRoundRobinPointer += 1
			continue
		}
		utils.LogWrite(self.logLabel, utils.LOG_DEBUG, fmt.Sprintf("FOUND AVAILABLE AGENT %d [serving %d customers]", pickedAgent.ID, pickedAgent.CustomerCount))
		self.agents[self.lastRoundRobinPointer].CustomerCount += 1
		if self.qiscusApi.AssignAgentToRoom(roomId, pickedAgent.ID) != nil {
			err = errors.New("FAILED TO ALLOCATE AGENT TO ROOM")
			break
		}

		utils.LogWrite(self.logLabel, utils.LOG_INFO, fmt.Sprintf("Agent %d allocated to room: %d", pickedAgent.ID, roomId))
		break
	}

	self.lastRoundRobinPointer += 1

	self.lock.Unlock()
	return
}
