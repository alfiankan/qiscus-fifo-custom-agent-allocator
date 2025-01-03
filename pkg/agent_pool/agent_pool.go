package agentpool

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
)

type AgentPoolConfig struct {
	MaxServedCustomerPerAgent int
	SyncInterval              int
	QiscusBaseHttpApiHost     string
	QiscusApiAuthAppId        string
	QiscusApiAuthSecret       string
}

type AgentPool struct {
	agents                []int
	config                AgentPoolConfig
	logLabel              string
	lock                  sync.RWMutex
	lastRoundRobinPointer int
	QiscusApi             QiscusApiClientInterface
}

// Create new agent allocator using roundrobin, this pool will keep updating at n interval time
// to qiscuss get all agent list as background sync
func NewAgentPoolAllocator(cfg AgentPoolConfig, qiscusApiClient QiscusApiClientInterface) *AgentPool {
	pool := &AgentPool{
		config:                cfg,
		agents:                []int{},
		logLabel:              "AGENT_ALLOCATOR",
		lastRoundRobinPointer: 0,
		QiscusApi:             qiscusApiClient,
	}
	pool.syncAgent()
	go pool.startBackgroundSync()

	return pool
}

func (self *AgentPool) GetTotalAgent() int {
	return len(self.agents)
}

// Ticker background to update agents list in the pool
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

	syncedAgent := []int{}

	pageNow := 1

	for {
		agentsResponseData, err := self.QiscusApi.GetAllAgents(pageNow, 10)
		if err != nil {
			break
		}
		if len(agentsResponseData.Data.Agents.Data) == 0 {
			break
		}

		for _, agent := range agentsResponseData.Data.Agents.Data {
			syncedAgent = append(syncedAgent, agent.ID)
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

// Allocate to available agent
func (self *AgentPool) AllocateAgent(roomId int) (err error) {
	utils.LogWrite(self.logLabel, utils.LOG_INFO, fmt.Sprintf("Allocating to room: %d", roomId))

	self.lock.Lock()

	// reset round robin if more than offset
	if len(self.agents) > self.lastRoundRobinPointer {
		self.lastRoundRobinPointer = 0
	}
	for {
		if self.lastRoundRobinPointer > len(self.agents)-1 {
			err = errors.New("ALL AGENTS ARE BUSY")
			self.lastRoundRobinPointer = 0
			self.lock.Unlock()
			return
		}
		pickedAgent := self.agents[self.lastRoundRobinPointer]

		// get latest customer served count from qiscusApi or continue to next agent ids in pool
		qiscusAgent, err := self.QiscusApi.GetAgentDetailById(pickedAgent)
		if err != nil {
			self.lastRoundRobinPointer += 1
			continue
		}

		if qiscusAgent.CurrentCustomerCount == self.config.MaxServedCustomerPerAgent || !qiscusAgent.IsAvailable {
			self.lastRoundRobinPointer += 1
			continue
		}
		utils.LogWrite(self.logLabel, utils.LOG_DEBUG, fmt.Sprintf("FOUND AVAILABLE AGENT %d [serving %d customers]", pickedAgent, qiscusAgent.CurrentCustomerCount))
		if self.QiscusApi.AssignAgentToRoom(roomId, pickedAgent) != nil {
			err = errors.New("FAILED TO ALLOCATE AGENT TO ROOM")
			break
		}

		utils.LogWrite(self.logLabel, utils.LOG_INFO, fmt.Sprintf("Agent %d allocated to room: %d", pickedAgent, roomId))
		break
	}

	self.lastRoundRobinPointer += 1
	self.lock.Unlock()
	return
}
