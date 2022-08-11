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
	path           string
	conn           *websocket.Conn
	wsModule       *WebsocketModule
	messageHandler func(message *MessageEvent, session *Session)
	closeHandler   func(closeMessage *CloseEvent, session *Session) error
	replyChannel   chan *ReplyMessage
}

func (session *Session) init() {
	session.conn.SetCloseHandler(func(code int, text string) error {
		closeMessage := makeCloseEvent(session, code, text)
		session.wsModule.events <- closeMessage
		return nil
	})
}

func makeSession(path string, conn *websocket.Conn, wsModule *WebsocketModule, messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) *Session {
	return &Session{path: path, conn: conn, wsModule: wsModule, messageHandler: messageHandler, closeHandler: closeHandler, replyChannel: make(chan *ReplyMessage)}
}

type MessageEvent struct {
	session     *Session
	MessageType int
	Data        []byte
	Err         error
}

func (messageEvent *MessageEvent) process() {
	if messageEvent.session != nil && messageEvent.session.messageHandler != nil {
		messageEvent.session.messageHandler(messageEvent, messageEvent.session)
	}
}

func makeMessageEvent(session *Session, messageType int, data []byte, err error) *MessageEvent {
	return &MessageEvent{session, messageType, data, err}
}

type CloseEvent struct {
	session *Session
	Code    int
	Text    string
}

func (closeEvent *CloseEvent) process() {
	if closeEvent.session != nil && closeEvent.session.closeHandler != nil {
		closeEvent.session.closeHandler(closeEvent, closeEvent.session)
	}
}

func makeCloseEvent(session *Session, code int, text string) *CloseEvent {
	return &CloseEvent{session, code, text}
}

type ReplyMessage struct {
	messageType int
	data        []byte
}

func MakeReplyMessage(messageType int, data []byte) *ReplyMessage {
	return &ReplyMessage{messageType, data}
}

type IOErrorEvent struct {
}

func (session *Session) listen() {
	for {
		messageType, p, err := session.conn.ReadMessage()
		message := makeMessageEvent(session, messageType, p, err)
		if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
			session.wsModule.events <- message
		}
		if err != nil {
			fmt.Println("[WS] Receive data from client failed", err)
			break
		}
	}
	fmt.Println("End ws read thread")
}

func (session *Session) response() {
	for replyMessage := range session.replyChannel {
		err := session.conn.WriteMessage(replyMessage.messageType, replyMessage.data)
		if err != nil {
			fmt.Println("[WS] Send data to client failed", err)
			break
		}
	}
	fmt.Println("End ws write thread")
}

func (session *Session) WriteText(data []byte) {
	replyMessage := MakeReplyMessage(websocket.TextMessage, data)
	session.replyChannel <- replyMessage
}

func (session *Session) WriteBytes(data []byte) {
	replyMessage := MakeReplyMessage(websocket.BinaryMessage, data)
	session.replyChannel <- replyMessage
}

func (session *Session) Close(code int, data string) {
	replyMessage := MakeReplyMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, data))
	session.replyChannel <- replyMessage
	close(session.replyChannel)

}

type WebsocketModule struct {
	BaseModule
	server   *http.Server
	sessions map[*Session]struct{}
}

func (websocket *WebsocketModule) makeWSHandler(path string, messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) {
	websocket.server.Handler.(*http.ServeMux).HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request come")
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		// conn.SetReadDeadline(time.Time{}.Add(time.Second))
		if err != nil {
			log.Println(err)
		} else {
			session := makeSession(path, conn, websocket, messageHandler, closeHandler)
			session.init()
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
