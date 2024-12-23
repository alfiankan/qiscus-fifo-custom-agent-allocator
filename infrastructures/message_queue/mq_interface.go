package messagequeue

import "time"

type handleNewChatQueued func(int) error

type MessageQueue interface {
	Push(roomId int) (err error)
	Pull(fn handleNewChatQueued, backOffInterval time.Duration, backOffMaxFail int) (err error)
}
