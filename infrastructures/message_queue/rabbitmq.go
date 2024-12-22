package messagequeue

import (
	"strconv"
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

func (self *MessageQueuRabbitMQ) Pull() (roomId int, err error) {
  return
}
