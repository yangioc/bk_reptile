package msg_nats

import (
	"context"

	"github.com/nats-io/nats.go"
)

func New(ctx context.Context, writeFlow chan []byte) *Handler {
	hdr := &Handler{
		ctx:       ctx,
		writeFlow: writeFlow,
		subMap:    make(map[string]*nats.Subscription),
		// subFlow:       make(chan *flow.Flow, constant.ChanSize_Large),
		// config:        config,
		// status:        status_init,
		// subjectMap:    sync.Map{},
		// writeEnd:      make(chan struct{}),
		// shutdownChan:  shutdownChan,
		// subscriberMap: make(map[string]func() error),
	}

	err := hdr.connect("nats://192.168.1.109:4223")
	if err != nil {
		return nil
	}

	// ch1 := make(chan *nats.Msg, 1024)
	// ch2 := make(chan *nats.Msg, 1024)
	// ch3 := make(chan *nats.Msg, 1024)

	// hdr.Sub("hello1", func(msg *nats.Msg) {
	// 	fmt.Printf("Sub Reply: %s\n", msg.Data)
	// })
	// hdr.SubChan("hello1", ch1)
	// go handlech("SubChan", ch1)

	// hdr.SubGroup("hello1", "h1", func(msg *nats.Msg) {
	// 	fmt.Printf("SubGroup Reply: %s\n", msg.Data)
	// })
	// hdr.SubGroupChan("hello1", "h2", ch2)
	// go handlech("SubGroupChan", ch2)

	// go hdr.SubSync("hello1")
	// go hdr.SubGroupSync("hello1", "h1", func(msg *nats.Msg) {
	// 	fmt.Printf("SubGroupSync Reply: %s\n", msg.Data)
	// })
	// go hdr.SubGroupChanSync("hello1", "h2", ch3)
	// go handlech("SubGroupChanSync", ch3)
	return hdr
}

// func cb(msg *nats.Msg) {

// 	// Use the response
// 	fmt.Printf("Reply: %s\n", msg.Data)
// }

// func handlech(tag string, ch chan *nats.Msg) {

// 	for {
// 		select {
// 		case msg := <-ch:
// 			// Use the response
// 			fmt.Printf("%s Reply: %s\n", tag, msg.Data)
// 		}
// 	}
// }
