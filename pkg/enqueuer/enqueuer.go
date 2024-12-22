package enqueuer

import (
	"fmt"
	"net/http"
	"strconv"

	messagequeue "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/infrastructures/message_queue"
	agentpool "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/pkg/agent_pool"
	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
	"github.com/labstack/echo/v4"
)

type WebHookEnqueuer struct {
	port                  int
	secret                string
	logLabel              string
	messageQueue          messagequeue.MessageQueue
	appId                 string
	appSecret             string
	maxCustServer         int
	agentPoolIntervalSync int
}

func NewWebHookHandlerEnqueuer(
	httpPort int,
	secret string,
	rabbitMqDSN string,
	appId string,
	appSecret string,
	maxCustServer int,
	agentPoolIntervalSync int,
) *WebHookEnqueuer {
	mq, err := messagequeue.NewMessageQueuRabbitMQ(rabbitMqDSN, "custom_allocator")
	if err != nil {
		utils.LogWrite("WEBHOOK_HTTP_ENQUEUER", utils.LOG_ERROR, "DIAlING AMQP", err.Error())
		panic(err)
	}
	return &WebHookEnqueuer{
		port:                  httpPort,
		secret:                secret,
		logLabel:              "WEBHOOK_HTTP_ENQUEUER",
		messageQueue:          mq,
		appId:                 appId,
		appSecret:             appSecret,
		maxCustServer:         maxCustServer,
		agentPoolIntervalSync: agentPoolIntervalSync,
	}
}

func (self *WebHookEnqueuer) RunWebhookHandler() {

	e := echo.New()
	e.POST("/allocate", func(c echo.Context) error {

		secret := c.QueryParam("secret")
		if secret != self.secret {
			return c.String(http.StatusUnauthorized, "UNAUTHORIZED")
		}

		var newMesage QiscusWebhookChatReqBody
		err := c.Bind(&newMesage)
		if err != nil {
			utils.LogWrite(self.logLabel, utils.LOG_ERROR, "WEBHOOK PARSE", err.Error())
			return c.String(http.StatusBadRequest, "Cant Parse Request Body")
		}
		fmt.Println(newMesage)

		roomId, err := strconv.Atoi(newMesage.RoomID)
		if err != nil {
			utils.LogWrite(self.logLabel, utils.LOG_ERROR, "WEBHOOK PARSE", err.Error())
			return c.String(http.StatusBadRequest, "Cant Parse Request Body")
		}

		if err := self.messageQueue.Push(roomId); err != nil {
			utils.LogWrite(self.logLabel, utils.LOG_ERROR, "ALLOCATOR ERROR", err.Error())
			return c.String(http.StatusInternalServerError, "Cant Enqueu Chat")
		}

		return c.String(http.StatusOK, "OK")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", self.port)))

}

func (self *WebHookEnqueuer) RunAllocator() {
	agentPoolCfg := agentpool.AgentPoolConfig{
		MaxServedCustomerPerAgent: self.maxCustServer,
		SyncInterval:              self.agentPoolIntervalSync,
		QiscusApiAuthAppId:        self.appId,
		QiscusApiAuthSecret:       self.appSecret,
		QiscusBaseHttpApiHost:     "https://multichannel.qiscus.com",
	}

	agentPool := agentpool.NewAgentPoolAllocator(agentPoolCfg)
	self.messageQueue.Pull(agentPool.AllocateAgent)

}
