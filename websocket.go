package heeglibs

import (
	"bytes"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var ClientConnects int

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	ClientChanSize = 32
	maxMessageSize = 1024 * 64

	RecvBufferSize = 64 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func init() {
	ClientConnects = 0
}

type Client struct {
	hub   *WebSocketServer
	conn  *websocket.Conn
	Send  chan []byte
	RCall func(*Client, []byte)
}

// websocket读取线程
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		// 将数据的\n用''来代替//
		// 同时处理数据包的大小 //
		// 定义数据格式 //
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.RCall(c, message)
		// 对消息进行处理
	}
}

// websocket连接写线程
func (this *Client) writePump() {
	defer func() {
		this.hub.unregister <- this
		this.conn.Close()
		log.Error("Client connect close!")
	}()

	for {
		select {
		case message, ok := <-this.Send:
			if !ok {
				log.Println("WebSocket writePump over? ok: ", ok)

				return
			}

			/**
			* 需要设置为发送二进制数据，否则浏览器会报出UTF-8解码错误
			 */
			w, err := this.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Println("WebSocket writePump over? NextWriter: ", err)

				return
			}

			w.Write(message)
			if err := w.Close(); err != nil {
				log.Println("WebSocket writePump over? Close: ", err)

				return
			}

			n := len(this.Send)
			for i := 0; i < n-1; i++ {
				data := <-this.Send
				w, err := this.conn.NextWriter(websocket.BinaryMessage)
				if err != nil {
					log.Println("WebSocket writePump over? NextWriter: ", err)

					continue
				}

				runtime.Gosched()
				w.Write(data)

				if err := w.Close(); err != nil {
					log.Println("WebSocket writePump over? Close: ", err)

					continue
				}
			}
		}
	}

}

// 创建websocket连接
// hub 		websocket服务对象
// w 		http请求的响应对象
// r 	 	http请求的请求对象
// rcall 	websocket读取数据的回调函数
// onCall 	websocket建立的回调函数
func WebSocketConn(hub *WebSocketServer, w http.ResponseWriter, r *http.Request, rcall func(*Client, []byte),
	onCall func(client *Client, state bool)) {
	var upgrade = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},

		ReadBufferSize:  1024 * 4,      // 4kb
		WriteBufferSize: 64 * 1024 * 1, // 64kb
	}

	/**
	* 设置响应头部，要不会报错，也就是不设置就连接不上
	 */
	var headers http.Header = make(http.Header)
	headers.Add("Sec-WebSocket-Protocol", "null")

	conn, err := upgrade.Upgrade(w, r, headers)
	if nil != err {
		log.Println(err)
		onCall(nil, false)

		return
	}

	client := &Client{hub: hub, conn: conn, Send: make(chan []byte, ClientChanSize), RCall: rcall}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()

	onCall(client, true)

	return
}

type WebSocketServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

// 创建性的WebsocketServer对象
func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// websocket建立的状态服务
// 主要是监听连接状态
func (h *WebSocketServer) RunServer() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			ClientConnects = ClientConnects + 1

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				log.Error("Client chan close!")

				ClientConnects = ClientConnects - 1
			}
		}
	}
}
