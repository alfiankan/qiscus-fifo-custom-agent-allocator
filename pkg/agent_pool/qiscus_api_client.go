package agentpool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
)

type QiscusApiClientInterface interface {
	GetAllAgents(page, perPage int) (agentsResponseData QiscusListAgentsApiResponse, err error)
	AssignAgentToRoom(roomId, agentId int) (err error)
	GetAgentDetailById(agentId int) (agent QiscusAget, err error)
}

type QiscusApiClient struct {
	httpClient *http.Client
	baseHost   string
	appId      string
	appSecret  string
	logLabel   string
}

func NewQiscusApiClient(baseHost string, appId string, appSecret string) QiscusApiClientInterface {
	return &QiscusApiClient{
		httpClient: &http.Client{},
		appId:      appId,
		appSecret:  appSecret,
		baseHost:   baseHost,
		logLabel:   "QISCUS_API_CLIENT",
	}
}

func (self *QiscusApiClient) GetAllAgents(page, perPage int) (agentsResponseData QiscusListAgentsApiResponse, err error) {

	uri := fmt.Sprintf("%s/api/v1/admin/agents?page=%d&limit=%d", self.baseHost, page, perPage)
	utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "CALL API", uri)
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS PREPARING API CALL ERROR", err.Error())
		return
	}

	request.Header.Set("Qiscus-App-Id", self.appId)
	request.Header.Set("Qiscus-Secret-Key", self.appSecret)

	response, err := self.httpClient.Do(request)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API CALL ERROR", err.Error())
		return
	}

	if response.StatusCode != 200 {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API CALL ERROR", response.Status)
		err = errors.New("Error Response API")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&agentsResponseData)
	if err != nil {
		body, _ := io.ReadAll(response.Body)
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API RESPONSE DECODE ERROR", err.Error(), response.Status, string(body))
		return
	}

	return

}

func (self *QiscusApiClient) AssignAgentToRoom(roomId, agentId int) (err error) {

	uri := fmt.Sprintf("%s/api/v1/admin/service/assign_agent", self.baseHost)
	utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "CALL API", uri)

	var param = url.Values{}
	param.Set("room_id", strconv.Itoa(roomId))
	param.Set("agent_id", strconv.Itoa(agentId))
	param.Set("max_agent", "1")
	var assignAgentPayload = bytes.NewBufferString(param.Encode())

	request, err := http.NewRequest("POST", uri, assignAgentPayload)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS PREPARING API CALL ERROR", err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Qiscus-App-Id", self.appId)
	request.Header.Set("Qiscus-Secret-Key", self.appSecret)

	response, err := self.httpClient.Do(request)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API CALL ERROR", err.Error())
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API CALL ERROR", string(body))
		err = errors.New("FAILED TO ASSIGN AGENT VIA API")
		return
	}
	utils.LogWrite(self.logLabel, utils.LOG_INFO, "SUCCESS TO ASSIGN AGENT")

	return

}

func (self *QiscusApiClient) GetAgentDetailById(agentId int) (agent QiscusAget, err error) {

	uri := fmt.Sprintf("%s/api/v1/admin/agents/get_by_ids?ids[]=%d", self.baseHost, agentId)
	utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "CALL API", uri)

	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS PREPARING API CALL ERROR", err.Error())
		return
	}
	request.Header.Set("Qiscus-App-Id", self.appId)
	request.Header.Set("Qiscus-Secret-Key", self.appSecret)

	response, err := self.httpClient.Do(request)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API CALL ERROR", err.Error())
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API CALL ERROR", string(body))
		err = errors.New("FAILED TO GET AGENT DEATAIL VIA API")
		return
	}
	utils.LogWrite(self.logLabel, utils.LOG_INFO, "SUCCESS TO GET AGENT DETAIL")

	var agentsResponseData QiscusAgentListDataApiResponse
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&agentsResponseData)
	if err != nil {
		body, _ := io.ReadAll(response.Body)
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "QISCUS API RESPONSE DECODE ERROR", err.Error(), response.Status, string(body))
		return
	}
	if len(agentsResponseData.Data) == 0 {
		err = errors.New("AGENT NOT FOUND")
		return
	}

	agent = agentsResponseData.Data[0]

	return

}
