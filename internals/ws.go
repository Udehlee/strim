package internals

import (
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsConn struct {
	ID   string
	WS   *websocket.Conn
	send chan []byte
	hub  *Hub
}

func NewWsConn(h *Hub) *WsConn {
	return &WsConn{
		ID:   uuid.New().String(),
		send: make(chan []byte),
		hub:  h,
	}
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFiles("template/index.html"))
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (wc *WsConn) Route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", wc.Index)
	mux.HandleFunc("/ws", wc.HandleWSConnection)
	return mux
}

func (wc *WsConn) Index(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Println("template didn't execute", nil)
	}

}

func (wc *WsConn) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("failed to upgrade connection")
		return
	}
	wc.WS = conn
	wc.hub.register <- wc

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
// and unregisters the client if connection is close
func (wc *WsConn) readPump() {
	defer func() {
		wc.hub.unregister <- wc
		wc.WS.Close()

	}()

	for {
		_, msg, err := wc.WS.ReadMessage()
		if err != nil {
			log.Println("failed to read message:", err)
			break
		}
		wc.hub.broadcast <- msg
	}

}
