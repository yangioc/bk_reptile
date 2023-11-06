package websocket

import (
	"bk_reptile/model/websocket/socketclient"
	"context"
	"fmt"

	"nhooyr.io/websocket"
)

type Client struct {
	socket   *socketclient.Handler
	callback socketclient.SocketManagerCallBack
}

func NewClient(callback socketclient.SocketManagerCallBack) *Client {
	return &Client{
		callback: callback,
	}
}

func (socket *Client) Launch(addr string) error {
	if len(addr) == 0 {
		return fmt.Errorf("[Websocket][Launch] addr Error addr: %v", addr)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	conn, _, err := websocket.Dial(ctx, addr, nil)
	if err != nil {
		return fmt.Errorf("[Websocket][Launch] Listen Error addr: %v", addr)
	}

	socket.socket = socketclient.New(ctx, conn, socket.callback)
	return socket.socket.Listen()
}

func (socket *Client) Send(ctx context.Context, message []byte) error {
	return socket.socket.Send(ctx, message)
}
