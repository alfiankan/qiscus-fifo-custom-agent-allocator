package messagequeue

import (
	"strconv"
	"time"

	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueuRabbitMQ struct {
	logLabel  string
	channel   *amqp.Channel
	queueName string
}

func NewMessageQueuRabbitMQ(amqpDsn string, queueName string) (MessageQueue, error) {

	mq := &MessageQueuRabbitMQ{
		logLabel:  "RABBIT_MQ_CLIENT",
		queueName: queueName,
	}

	conn, err := amqp.Dial(amqpDsn)
	if err != nil {
		utils.LogWrite(mq.logLabel, utils.LOG_ERROR, "DIAL TO RABBITMQ", err.Error())
		return nil, err
	}

	mq.channel, err = conn.Channel()

	_, err = mq.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		utils.LogWrite(mq.logLabel, utils.LOG_ERROR, "DECLARING QUEUE", err.Error())
		return nil, err
	}

	return mq, nil
}

func (self *MessageQueuRabbitMQ) Push(roomId int) (err error) {
	err = self.channel.Publish(
		"",
		self.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.Itoa(roomId)),
		})
	return
}

func (self *MessageQueuRabbitMQ) Pull(fn handleNewChatQueued) (err error) {

	msgs, err := self.channel.Consume(
		self.queueName,
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

  deepSleepMaxTried := 100
  deepSleepThreshold := 0
  deepSleepInterval := 10 * time.Second

	for msg := range msgs {
    if deepSleepThreshold > deepSleepMaxTried {
      time.Sleep(deepSleepInterval)
      deepSleepThreshold = 0
    }

		utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "RECEIVING QUEUE", string(msg.Body))
		roomId, err := strconv.Atoi(string(msg.Body))
		if err != nil {
			utils.LogWrite(self.logLabel, utils.LOG_ERROR, "CANT PARSE QUEUE", string(msg.Body))
			if err = msg.Nack(false, true); err != nil {
				utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "UNACK", err.Error())
			}
		}
		if err := fn(roomId); err != nil {
      deepSleepThreshold += 1
			utils.LogWrite(self.logLabel, utils.LOG_ERROR, "CANT ALLOCATE AGENT - UNACK MESSAGE", err.Error())
			if err = msg.Nack(false, true); err != nil {
				utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "UNACK", err.Error())
			}
		} else {

			if err = msg.Ack(false); err != nil {
				utils.LogWrite(self.logLabel, utils.LOG_DEBUG, "ACK", err.Error())
			}
		}
	}

	return
}
