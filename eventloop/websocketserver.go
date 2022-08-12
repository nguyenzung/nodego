package eventloop

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  WS_READ_BUFFER_SIZE,
	WriteBufferSize: WS_WRITE_BUFFER_SIZE,
}

type Session struct {
	path           string
	conn           *websocket.Conn
	wsModule       *WebsocketModule
	messageHandler func(message *MessageEvent, session *Session)
	closeHandler   func(closeMessage *CloseEvent, session *Session) error
	replyChannel   chan *ReplyMessage
	isClosed       bool
	closeMutex     sync.Mutex
}

func (session *Session) init() {
	session.conn.SetCloseHandler(func(code int, text string) error {
		closeMessage := makeCloseEvent(session, code, text)
		session.wsModule.events <- closeMessage
		return nil
	})
}

func makeSession(path string, conn *websocket.Conn, wsModule *WebsocketModule, messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) *Session {
	return &Session{path: path, conn: conn, wsModule: wsModule, messageHandler: messageHandler, closeHandler: closeHandler, replyChannel: make(chan *ReplyMessage), isClosed: false}
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

func (session *Session) listen() {
	for {
		messageType, p, err := session.conn.ReadMessage()
		message := makeMessageEvent(session, messageType, p, err)
		session.wsModule.events <- message
		if err != nil {
			session.WriteClose(1000, "")
			break
		}
	}
}

func (session *Session) response() {
	for replyMessage := range session.replyChannel {
		err := session.conn.WriteMessage(replyMessage.messageType, replyMessage.data)
		if err != nil {
			fmt.Println("[WS SEND ERRO]", err)
		}
	}
}

func (session *Session) send(replyMessage *ReplyMessage) {
	session.closeMutex.Lock()
	defer session.closeMutex.Unlock()
	if !session.isClosed {
		session.replyChannel <- replyMessage
		if replyMessage.messageType == websocket.CloseMessage {
			session.isClosed = true
			close(session.replyChannel)
		}
	} else {
		fmt.Println("[WS ]Writing channel was closed")
	}
}

func (session *Session) WriteText(data []byte) {
	replyMessage := MakeReplyMessage(websocket.TextMessage, data)
	session.send(replyMessage)
}

func (session *Session) WriteBytes(data []byte) {
	replyMessage := MakeReplyMessage(websocket.BinaryMessage, data)
	session.send(replyMessage)
}

func (session *Session) WriteClose(code int, data string) {
	replyMessage := MakeReplyMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, data))
	session.send(replyMessage)
}

type WebsocketModule struct {
	BaseModule
	server     *http.Server
	numHandler int
}

func (websocket *WebsocketModule) makeWSHandler(path string, messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) {
	websocket.numHandler++
	websocket.server.Handler.(*http.ServeMux).HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		conn.SetReadDeadline(time.Now().Add(time.Second * WS_TIMEOUT_IN_SECONDS))
		if err != nil {
			log.Println(err)
		} else {
			session := makeSession(path, conn, websocket, messageHandler, closeHandler)
			session.init()
			go session.listen()
			go session.response()
		}
		log.Println("[WS] Client Connected")
	})
}

func (websocket *WebsocketModule) exec() {
	if websocket.numHandler > 0 {
		log.Fatal(websocket.server.ListenAndServe())
	}
}

func makeWebsocketModule(events chan IEvent) *WebsocketModule {
	server := makeServer(fmt.Sprintf("%s:%d", WEBSOCKET_IP, WEBSOCKET_PORT))
	return &WebsocketModule{BaseModule: BaseModule{events}, server: server}
}
