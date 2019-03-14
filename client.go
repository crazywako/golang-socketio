package gosocketio

import (
	"github.com/crazywako/golang-socketio/transport"
	"github.com/crazywako/golang-socketio/protocol"
	"fmt"
	"strconv"
)

const (
	webSocketProtocol = "ws://"
	webSocketSecureProtocol = "wss://"
	socketioUrl       = "/socket.io/?EIO=3&transport=websocket"
)

/**
Socket.io client representation
*/
type Client struct {
	methods
	Channel
}

/**
Get ws/wss url by host and port
 */
func GetUrl(host string, port int, secure bool) string {
	var prefix string
	if secure {
		prefix = webSocketSecureProtocol
	} else {
		prefix = webSocketProtocol
	}
	return prefix + host + ":" + strconv.Itoa(port) + socketioUrl
}

/**
connect to host and initialise socket.io protocol

The correct ws protocol url example:
ws://myserver.com/socket.io/?EIO=3&transport=websocket

You can use GetUrlByHost for generating correct url
*/
func Dial(url string, nsp string, tr transport.Transport) (*Client, error) {
	c := &Client{}
	c.initChannel()
	c.initMethods()

	var err error

	c.conn, err = tr.Connect(url)
	nspMsg := fmt.Sprintf("4%d%s", protocol.MessageTypeOpen, nsp)
    	c.conn.WriteMessage(nspMsg)
	if err != nil {
		return nil, err
	}

	go inLoop(&c.Channel, &c.methods)
	go outLoop(&c.Channel, &c.methods)
	go pinger(&c.Channel)

	return c, nil
}

/**
Close client connection
*/
func (c *Client) Close() {
	closeChannel(&c.Channel, &c.methods)
}
