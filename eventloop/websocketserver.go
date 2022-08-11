package eventloop

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Session struct {
	path             string
	conn             *websocket.Conn
	messageHandler   func(message MessageEvent, session *Session)
	closeHandler     func(code int, text string) error
	replyChannel     chan *ReplyMessage
	onMessageChannel chan *MessageEvent
}

func makeSession(path string, conn *websocket.Conn, messageHandler func(MessageEvent, *Session), closeHandler func(code int, text string) error) *Session {
	return &Session{path: path, conn: conn, messageHandler: messageHandler, closeHandler: closeHandler, replyChannel: make(chan *ReplyMessage), onMessageChannel: make(chan *MessageEvent)}
}

type MessageEvent struct {
	messageType int
	data        []byte
	err         error
}

func MakeMessageEvent(messageType int, data []byte, err error) *MessageEvent {
	return &MessageEvent{messageType, data, err}
}

type ReplyMessage struct {
	messageType int
	data        []byte
}

func MakeReplyMessage(messageType int, data []byte) *ReplyMessage {
	return &ReplyMessage{messageType, data}
}

func (session *Session) listen() {
	for {
		messageType, p, err := session.conn.ReadMessage()
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println("[Message]", "%d %s", messageType, string(p))
			message := MakeMessageEvent(messageType, p, err)
			session.onMessageChannel <- message
		}
	}
}

func (session *Session) response() {
	for replyMessage := range session.replyChannel {
		session.conn.WriteMessage(replyMessage.messageType, replyMessage.data)
	}
}

func (session *Session) WriteText(data []byte) {
	replyMessage := MakeReplyMessage(websocket.TextMessage, data)
	session.replyChannel <- replyMessage
}

func (session *Session) WriteBytes(data []byte) {
	replyMessage := MakeReplyMessage(websocket.BinaryMessage, data)
	session.replyChannel <- replyMessage
}

func (session *Session) Close(code int, data []byte) {
	replyMessage := MakeReplyMessage(websocket.BinaryMessage, data)
	session.replyChannel <- replyMessage
	close(session.replyChannel)
}

type WebsocketModule struct {
	BaseModule
	server   *http.Server
	sessions map[*Session]struct{}
}

func (websocket *WebsocketModule) MakeWS(path string, messageHandler func(message MessageEvent, session *Session), closeHandler func(code int, text string) error) {
	websocket.server.Handler.(*http.ServeMux).HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request come")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
		} else {
			session := makeSession(path, conn, messageHandler, closeHandler)
			go session.listen()
			go session.response()
		}

		log.Println("Client Connected")
	})
}

func (websocket *WebsocketModule) exec() {
	log.Fatal(websocket.server.ListenAndServe())
}

func makeWebsocketModule(events chan IEvent) *WebsocketModule {
	server := makeServer(fmt.Sprintf("%s:%d", WEBSOCKET_IP, WEBSOCKET_PORT))
	return &WebsocketModule{BaseModule: BaseModule{events}, server: server, sessions: make(map[*Session]struct{})}
}
