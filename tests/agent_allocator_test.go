package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"

	agentpool "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/pkg/agent_pool"
)

func TestRoundRobinAllocator(t *testing.T) {
	qiscusApiMock := agentpool.QiscusApiClientMock{}
	agents := agentpool.QiscusListAgentsApiResponse{}
	agents.Data.Agents.Data = []agentpool.QiscusAget{}

	agents.Data.Agents.Data = append(agents.Data.Agents.Data, agentpool.QiscusAget{
		ID:                   10,
		CurrentCustomerCount: 0,
		IsAvailable:          true,
	})
	qiscusApiMock.GetAllAgentsRet = &agents

	qiscusApiMock.GetAgentDetailByIdRet = &agents.Data.Agents.Data[0]

	agentAllocator := agentpool.NewAgentPoolAllocator(agentpool.AgentPoolConfig{
		MaxServedCustomerPerAgent: 2,
		SyncInterval:              500,
	}, &qiscusApiMock)

	agentAllocator.QiscusApi = &qiscusApiMock

	assert.Nil(t, agentAllocator.AllocateAgent(1))
	agents.Data.Agents.Data[0].CurrentCustomerCount = 1
	assert.Error(t, agentAllocator.AllocateAgent(2))
	assert.Nil(t, agentAllocator.AllocateAgent(2))

}
