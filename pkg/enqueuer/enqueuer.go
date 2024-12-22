package enqueuer

import (
	"fmt"
	"net/http"
	"strconv"

	messagequeue "github.com/alfiankan/qiscus-fifo-custom-agent-allocator/infrastructures/message_queue"
	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
	"github.com/labstack/echo/v4"
)

type WebHookEnqueuer struct {
	port         int
	secret       string
	logLabel     string
	messageQueue messagequeue.MessageQueue
}

func NewWebHookHandlerEnqueuer(httpPort int, secret string, rabbitMqDSN string) *WebHookEnqueuer {
  mq, err := messagequeue.NewMessageQueuRabbitMQ(rabbitMqDSN, "custom_allocator")
  if err != nil {
    utils.LogWrite("WEBHOOK_HTTP_ENQUEUER", utils.LOG_ERROR, "DIAlING AMQP", err.Error())
    panic(err)
  }
	return &WebHookEnqueuer{
		port:     httpPort,
		secret:   secret,
		logLabel: "WEBHOOK_HTTP_ENQUEUER",
    messageQueue: mq,
	}
}

func (self *WebHookEnqueuer) Run() {

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
