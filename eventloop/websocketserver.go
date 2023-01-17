package eventloop

import (
	"fmt"
	"log"
	"net/http"
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
	openHandler    func(session *Session)
	closeHandler   func(closeMessage *CloseEvent, session *Session) error
	replyChannel   chan *ReplyMessage
	isClosed       bool
}

func (session *Session) init() {
	go session.listen()
	go session.response()
	session.conn.SetCloseHandler(func(code int, text string) error {
		closeMessage := makeCloseEvent(session, code, text)
		session.wsModule.events <- closeMessage
		return nil
	})
}

func makeSession(path string, conn *websocket.Conn, wsModule *WebsocketModule, openHandler func(*Session), messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) *Session {
	return &Session{path: path, conn: conn, wsModule: wsModule, openHandler: openHandler, messageHandler: messageHandler, closeHandler: closeHandler, replyChannel: make(chan *ReplyMessage), isClosed: false}
}

type OpenEvent struct {
	session *Session
}

func (openEvent *OpenEvent) process() {
	if openEvent.session != nil && openEvent.session.openHandler != nil {
		openEvent.session.openHandler(openEvent.session)
	}
}

func makeOpenEvent(session *Session) *OpenEvent {
	return &OpenEvent{session}
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
	for !session.isClosed {
		messageType, p, err := session.conn.ReadMessage()
		log.Println(" ----->", messageType, p, err)
		switch messageType {
		case websocket.TextMessage:
			if err == nil {
				message := makeMessageEvent(session, messageType, p, err)
				session.wsModule.events <- message
			}

		case websocket.BinaryMessage:
			{

			}

		case websocket.PingMessage:
			{

			}

		case websocket.PongMessage:
			{

			}

		case websocket.CloseMessage:
			{

				session.handleCloseRequest()
			}

		default:
			{
				session.onTerminateSession()
			}
		}
	}
	log.Println("Read is closed")
}

func (session *Session) response() {
	for replyMessage := range session.replyChannel {
		err := session.conn.WriteMessage(replyMessage.messageType, replyMessage.data)
		if err != nil {
			fmt.Println("[WS SEND ERRO]", err)
		}
	}
	log.Println("Write is closed")
}

func (session *Session) send(replyMessage *ReplyMessage) {
	if !session.isClosed {
		session.replyChannel <- replyMessage
	} else {
		log.Println("[WS ]Writing channel was closed")
	}
}

func (session *Session) CloseSession(code int, data string) {
	session.WriteClose(code, data)
}

func (session *Session) handleCloseRequest() {
	log.Println("Handle Close request")
	session.onCloseSession()
}

func (session *Session) onCloseSession() {
	log.Println("on Close")
	session.isClosed = true
	close(session.replyChannel)
	// session.conn.Close()
}

func (session *Session) onTerminateSession() {
	log.Println("on Terminate")
	session.isClosed = true
	close(session.replyChannel)
	// session.conn.Close()
}

func (session *Session) WriteText(data []byte) {
	replyMessage := MakeReplyMessage(websocket.TextMessage, data)
	go session.send(replyMessage)
}

func (session *Session) WriteBytes(data []byte) {
	replyMessage := MakeReplyMessage(websocket.BinaryMessage, data)
	go session.send(replyMessage)
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

func (websocket *WebsocketModule) makeWSHandler(path string, openHandler func(*Session), messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) {
	websocket.numHandler++
	websocket.server.Handler.(*http.ServeMux).HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		conn.SetReadDeadline(time.Now().Add(time.Second * 20))
		if err != nil {
			log.Println(err)
		} else {
			session := makeSession(path, conn, websocket, openHandler, messageHandler, closeHandler)
			openEvent := makeOpenEvent(session)
			session.wsModule.events <- openEvent
			session.init()
		}
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
