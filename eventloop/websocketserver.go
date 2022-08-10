package eventloop

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Session struct {
	path         string
	conn         *websocket.Conn
	running      bool
	runningMutex sync.Mutex
	handler      func(message string, session *Session)
}

func makeSession(path string, conn *websocket.Conn, handler func(string, *Session)) *Session {
	return &Session{path, conn, true, sync.Mutex{}, handler}
}

func (session *Session) listen() {
	for session.isRunning() {
		messageType, p, err := session.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("[Message]", "%d %s", messageType, string(p))
	}
	session.conn.Close()
}

func (session *Session) isRunning() bool {
	session.runningMutex.Lock()
	defer session.runningMutex.Unlock()
	return session.running
}

func (session *Session) Close() {
	session.runningMutex.Lock()
	session.running = false
	session.runningMutex.Unlock()
}

func (session *Session) WriteText(data []byte) {
	session.conn.WriteMessage(websocket.TextMessage, data)
}

func (session *Session) WriteBytes(data []byte) {
	session.conn.WriteMessage(websocket.BinaryMessage, data)
}

func MakeWebsocketSession(path string, conn *websocket.Conn) *Session {
	return &Session{path: path, conn: conn}
}

func MakeWS(path string, messageHandler func(message string, session *Session)) {
	websocketModule.server.Handler.(*http.ServeMux).HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request come")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)

		conn.SetPingHandler(func(appData string) error {
			return nil
		})

		conn.SetPongHandler(func(appData string) error {
			return nil
		})

		conn.SetCloseHandler(func(code int, text string) error {
			return nil
		})

		if err != nil {
			log.Println(err)
		} else {
			session := makeSession(path, conn, messageHandler)
			session.listen()
		}

		log.Println("Client Connected")
	})
}

type WebsocketModule struct {
	BaseModule
	server   *http.Server
	sessions map[*Session]struct{}
}

func (websocket *WebsocketModule) exec() {
	log.Fatal(websocket.server.ListenAndServe())
}

var websocketModule *WebsocketModule

func initWebsocketModule(events chan IEvent) {
	server := makeServer(fmt.Sprintf("%s:%d", WEBSOCKET_IP, WEBSOCKET_PORT))
	websocketModule = &WebsocketModule{BaseModule: BaseModule{events}, server: server, sessions: make(map[*Session]struct{})}
}

func startWebsocketModule() {
	go websocketModule.exec()
}
