package eventloop

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Session struct {
	path string
	conn *websocket.Conn
}

func (session *Session) WriteBytes(data []byte) {
	session.conn.WriteMessage(websocket.BinaryMessage, data)
}

func (session *Session) Close() {
	session.conn.Close()
}

func MakeWebsocketSession(path string, conn *websocket.Conn) *Session {
	return &Session{path: path, conn: conn}
}

func MakeWS(path string, handler func(message string, session *Session)) {
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {

	})
}

type WebsocketModule struct {
	BaseModule
	sessions map[*Session]struct{}
}

func (websocket *WebsocketModule) exec() {
	log.Fatal()
}

var websocketModule *WebsocketModule

func initWebsocketModule(events chan IEvent) {
	websocketModule = &WebsocketModule{BaseModule: BaseModule{events}, sessions: make(map[*Session]struct{})}
}

func startWebsocketModule() {
	go websocketModule.exec()
}
