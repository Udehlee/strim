package internals

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	WS   *websocket.Conn
	send chan []byte
}

func NewWsConn(wsConn *websocket.Conn) *WsConn {
	return &WsConn{
		WS:   wsConn,
		send: make(chan []byte),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wc *WsConn) Route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", wc.HandleWSConnection)
	return mux
}

func (wc *WsConn) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection")
		return
	}

	wc.WS = conn

	go func() {
		for msg := range wc.send {
			err := wc.WS.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("failed to write message", err)
				break
			}

		}

	}()

	wc.readPump()
}

// readPump reads incoming messages from the WebSocket connection
func (wc *WsConn) readPump() {
	for {
		_, _, err := wc.WS.ReadMessage()
		if err != nil {
			log.Println("failed to read message:", err)
			break
		}

	}
}
