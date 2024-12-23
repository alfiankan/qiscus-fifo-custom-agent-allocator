package agentpool

import "errors"

type QiscusApiClientMock struct {
	GetAllAgentsRet          *QiscusListAgentsApiResponse
	GetAllAgentsRetErr       error
	AssignAgentToRoomRetErr  error
	GetAgentDetailByIdRet    *QiscusAget
	GetAgentDetailByIdRetErr error
}

func NewQiscusApiClientMock() QiscusApiClientInterface {
	return &QiscusApiClientMock{}
}

func (self *QiscusApiClientMock) GetAllAgents(page, perPage int) (agentsResponseData QiscusListAgentsApiResponse, err error) {
	if page == 2 {
		err = errors.New("EOF")
		return
	}
	return *self.GetAllAgentsRet, self.GetAllAgentsRetErr
}

func (self *QiscusApiClientMock) AssignAgentToRoom(roomId, agentId int) (err error) {
	return self.GetAgentDetailByIdRetErr
}

func (self *QiscusApiClientMock) GetAgentDetailById(agentId int) (agent QiscusAget, err error) {
	return *self.GetAgentDetailByIdRet, self.GetAgentDetailByIdRetErr
}
