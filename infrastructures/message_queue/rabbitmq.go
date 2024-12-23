package messagequeue

import (
	"strconv"
	"time"

	"github.com/alfiankan/qiscus-fifo-custom-agent-allocator/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueuRabbitMQ struct {
	logLabel   string
	channel    *amqp.Channel
	connection *amqp.Connection
	queueName  string
	amqpDsn    string
}

func NewMessageQueuRabbitMQ(amqpDsn string, queueName string) (mq MessageQueue, err error) {

	mqClient := MessageQueuRabbitMQ{
		logLabel:  "RABBIT_MQ_CLIENT",
		queueName: queueName,
		amqpDsn:   amqpDsn,
	}

	mqClient.Connect()

	// reconnect
	go func() {

		for {
			time.Sleep(1 * time.Second)
			if mqClient.connection.IsClosed() {
				for {
					time.Sleep(10 * time.Second)
					utils.LogWrite(mqClient.logLabel, utils.LOG_INFO, "TRYING REDIAL TO RABBITMQ")

					err := mqClient.Connect()

					if err != nil {
						utils.LogWrite(mqClient.logLabel, utils.LOG_ERROR, "REDIAL ERROR TO RABBITMQ", err.Error())
					} else {
						utils.LogWrite(mqClient.logLabel, utils.LOG_INFO, "REDIAL SUCCESS TO RABBITMQ")
						break
					}
				}
			}
		}

	}()

	_, err = mqClient.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		utils.LogWrite(mqClient.logLabel, utils.LOG_ERROR, "DECLARING QUEUE", err.Error())
		return nil, err
	}

	return &mqClient, nil
}

func (self *MessageQueuRabbitMQ) Connect() (err error) {
	connection, err := amqp.Dial(self.amqpDsn)
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "DIAL TO RABBITMQ", err.Error())
		return
	}
	channel, err := connection.Channel()
	if err != nil {
		utils.LogWrite(self.logLabel, utils.LOG_ERROR, "DIAL TO RABBITMQ", err.Error())
		return
	}
	self.connection = connection
	self.channel = channel
	return
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

func (self *MessageQueuRabbitMQ) Pull(
	fn handleNewChatQueued,
	backOffInterval time.Duration,
	backOffMaxFail int,
) (err error) {

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

	deepSleepThreshold := 0

	for msg := range msgs {
		if deepSleepThreshold > backOffMaxFail {
			time.Sleep(time.Duration(backOffInterval))
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
