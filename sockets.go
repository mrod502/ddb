package ddb

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mrod502/logger"
)

//public vars
var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		CheckOrigin:       checkOrigin,
		EnableCompression: true,
	}
	StreamExtension = "/ws"
	send            chan []byte
	writeTimeout    = time.Second
	pongWait        = time.Minute
	PingWait        = time.Second * 30
	ErrUnsubscribe  = errors.New("unsubscribed")
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients        map[*Client]bool
	broadcast      chan []byte
	register       chan *Client
	unregister     chan *Client
	mux            sync.RWMutex
	readHandleFunc func([]byte) error
}

func newHub(f func([]byte) error) *Hub {
	return &Hub{
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		readHandleFunc: f,
	}
}

func wssServe(h *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Sockets", "Upgrade", err.Error())
	}
	client := &Client{hub: h, conn: conn}
	client.hub.register <- client
	go client.write()
}

func checkOrigin(r *http.Request) bool {
	return true
}

func (h *Hub) serve() {
	for {
		select {
		case msg := <-h.broadcast:
			h.mux.RLock()
			for client := range h.clients {
				client.send <- msg
			}
			h.mux.RUnlock()
		case client := <-h.register:
			h.mux.Lock()
			h.clients[client] = true
			h.mux.Unlock()
		case client := <-h.unregister:
			_ = client.conn.Close()
			h.mux.Lock()
			delete(h.clients, client)
			h.mux.Unlock()
		}
	}
}

func (c *Client) write() {

	for {
		select {
		case msg := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Warn("wssServer", fmt.Sprintf("%s error:", c.conn.UnderlyingConn().RemoteAddr()), err.Error())
				c.hub.unregister <- c
				return
			}
		default:
		}
	}
}

func (c *Client) read() {

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.conn.SetPingHandler(DefaultPingHandler)

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("wssServe", fmt.Sprintf("client:%s error:", c.conn.UnderlyingConn().RemoteAddr()), err.Error())
			}
			break
		}
		err = c.hub.readHandleFunc(msg)
		if err == ErrUnsubscribe {
			logger.Warn("WssServe", "unsubscribe", c.conn.UnderlyingConn().RemoteAddr().String(), err.Error())
			c.conn.Close()
			c.hub.unregister <- c
			return
		}
	}

}

type DBSubscription struct {
	Endpoints []string
}

func DefaultPingHandler(s string) error {

	return nil
}
