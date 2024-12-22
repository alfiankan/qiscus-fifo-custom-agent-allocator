package messagequeue

type MessageQueue interface {
	Push(roomId int) (err error)
	Pull() (roomId int, err error)
}
