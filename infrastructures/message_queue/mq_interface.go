package messagequeue

type handleNewChatQueued func(int) error

type MessageQueue interface {
	Push(roomId int) (err error)
	Pull(fn handleNewChatQueued) (err error)
}
