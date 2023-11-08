package schedule

type MessageQ interface{}

type Handle struct {
	mq MessageQ
}

func New(mq MessageQ) *Handle {
	return &Handle{mq: mq}
}
