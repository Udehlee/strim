package internals

type Hub struct {
	clients    map[string]*WsConn
	register   chan *WsConn
	unregister chan *WsConn
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*WsConn),
		register:   make(chan *WsConn),
		unregister: make(chan *WsConn),
		broadcast:  make(chan []byte),
	}
}

// registerConn adds a new client to active connections
func (h Hub) registerConn(conn *WsConn) {
	h.clients[conn.ID] = conn
}

// unregisterConn removes a client active connections
func (h *Hub) unregisterConn(conn *WsConn) {
	_, ok := h.clients[conn.ID]
	if ok {
		delete(h.clients, conn.ID)
		close(conn.send)
	}
}

// broadcastMsg sends message to all connected clients
// removes disconnected clients
func (h *Hub) broadcastMsg(msg []byte) {
	for _, client := range h.clients {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(h.clients, client.ID)
		}

	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.registerConn(conn)
		case conn := <-h.unregister:
			h.unregisterConn(conn)
		case msg := <-h.broadcast:
			h.broadcastMsg(msg)
		}
	}
}
