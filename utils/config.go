package utils

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	QiscusApiBaseHost               string
	AppId                           string
	AppSecret                       string
	WebHookSecret                   string
	WebhookPort                     int
	Amqp                            string
	MaxCustServer                   int
	AgentPoolIntervalSync           int
	QueueBackoffSleepIntervalSecond int
}

func LoadApplicationConfig() *ApplicationConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := ApplicationConfig{}
	cfg.AppId = os.Getenv("QISCUS_APP_ID")
	cfg.AppSecret = os.Getenv("QISCUS_APP_SECRET")
	cfg.WebHookSecret = os.Getenv("SECRET")
	cfg.WebhookPort, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.Amqp = os.Getenv("RABBIT_MQ_DSN")
	cfg.MaxCustServer, _ = strconv.Atoi(os.Getenv("MAX_SERVED_CUSTOMER_PER_AGENT"))
	cfg.AgentPoolIntervalSync, _ = strconv.Atoi(os.Getenv("AGENT_POOL_SYNC_INTERVAL"))
	cfg.QueueBackoffSleepIntervalSecond, _ = strconv.Atoi(os.Getenv("QUEUE_BACKOFF_SLEEP_INTERVAL_SECOND"))
	cfg.QiscusApiBaseHost = os.Getenv("QISCUS_API_HOST")
	return &cfg

}
