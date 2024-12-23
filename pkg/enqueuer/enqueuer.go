package enqueuer

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	messagequeue "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/infrastructures/message_queue"
	agentpool "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/pkg/agent_pool"
	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
	"github.com/labstack/echo/v4"
)

type WebHookEnqueuer struct {
	messageQueue messagequeue.MessageQueue
	cfg          *utils.ApplicationConfig
	logLabel     string
}

func NewWebHookHandlerEnqueuer(cfg *utils.ApplicationConfig) *WebHookEnqueuer {

	mq, err := messagequeue.NewMessageQueuRabbitMQ(cfg.Amqp, "custom_allocator")
	if err != nil {
		utils.LogWrite("WEBHOOK_HTTP_ENQUEUER", utils.LOG_ERROR, "DIAlING AMQP", err.Error())
		panic(err)
	}
	return &WebHookEnqueuer{
		logLabel:     "WEBHOOK_HTTP_ENQUEUER",
		messageQueue: mq,
		cfg:          cfg,
	}
}

func (self *WebHookEnqueuer) RunWebhookHandler() {

	e := echo.New()
	e.POST("/allocate", func(c echo.Context) error {

		secret := c.QueryParam("secret")
		if secret != self.cfg.WebHookSecret {
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
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", self.cfg.WebhookPort)))

}

func (self *WebHookEnqueuer) RunAllocator() {
	agentPoolCfg := agentpool.AgentPoolConfig{
		MaxServedCustomerPerAgent: self.cfg.MaxCustServer,
		SyncInterval:              self.cfg.AgentPoolIntervalSync,
		QiscusApiAuthAppId:        self.cfg.AppId,
		QiscusApiAuthSecret:       self.cfg.AppSecret,
		QiscusBaseHttpApiHost:     self.cfg.QiscusApiBaseHost,
	}

	agentPool := agentpool.NewAgentPoolAllocator(agentPoolCfg)
	self.messageQueue.Pull(
		agentPool.AllocateAgent,
		time.Duration(self.cfg.QueueBackoffSleepIntervalSecond)*time.Second,
		agentPool.GetTotalAgent()*2,
	)
}
