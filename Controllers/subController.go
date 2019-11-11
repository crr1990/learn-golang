package Controllers
//
//import (
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//)
//
//type SSEHandler struct {
//	// Create a map of clients, the keys of the map are the channels over
//	// which we can push messages to attached clients. (The values are just
//	// booleans and are meaningless.)
//	clients map[chan string]bool
//
//	// Channel into which new clients can be pushed
//	newClients chan chan string
//
//	// Channel into which disconnected clients should be pushed
//	defunctClients chan chan string
//
//	// Channel into which messages are pushed to be broadcast out
//	messages chan string
//}
//
//// Make a new SSEHandler instance.
//func NewSSEHandler() *SSEHandler {
//	b := &SSEHandler{
//		clients:        make(map[chan string]bool),
//		newClients:     make(chan (chan string)),
//		defunctClients: make(chan (chan string)),
//		messages:       make(chan string, 10), // buffer 10 msgs and don't block sends
//	}
//
//	go b.HandleEvents()
//	return b
//}
//
//// Start handling new and disconnected clients, as well as sending messages to
//// all connected clients.
//func (b *SSEHandler) HandleEvents() {
//	go func() {
//		for {
//			select {
//			case s := <-b.newClients:
//				b.clients[s] = true
//				log.Printf("Client added. %d registered clients", len(b.clients))
//			case s := <-b.defunctClients:
//				delete(b.clients, s)
//				//close(s)
//				log.Printf("Client remove. %d registered clients", len(b.clients))
//			case msg := <-b.messages:
//				log.Printf("Client msg. %s ", msg)
//				for s, _ := range b.clients {
//					s <- msg
//				}
//			}
//		}
//	}()
//}
//
//// Send out a simple string to all clients.
//func (b *SSEHandler) SendString(msg string) {
//	b.messages <- msg
//}
//
//// Send out a JSON string object to all clients.
//func (b *SSEHandler) SendJSON(obj interface{}) {
//	tmp, err := json.Marshal(obj)
//	if err != nil {
//		log.Panic("Error while sending JSON object:", err)
//	}
//	b.messages <- string(tmp)
//}
//
//func Sub(c *gin.Context)  {
//	 b := NewSSEHandler()
//	 b.Subscribe(c)
//	 b.SendString("hi")
//}
//
//// Subscribe a new client and start sending out messages to it.
//func (b *SSEHandler) Subscribe(c *gin.Context) {
//	w := c.Writer
//	f, ok := w.(http.Flusher)
//	if !ok {
//		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Streaming unsupported"))
//		return
//	}
//
//	c.Writer.Header().Set("Content-Type","text/event-stream")
//	c.Writer.Header().Set("Cache-Control","no-cache")
//	c.Writer.Header().Set("Connection","keep-alive")
//	c.Writer.Header().Set("Access-Control-Allow-Origin","*")
//
//
//
//	// Create a new channel, over which we can send this client messages.
//	messageChan := make(chan string)
//	// Add this client to the map of those that should receive updates
//	b.newClients <- messageChan
//
//	//notify := w.(http.CloseNotifier).CloseNotify()
//	notify := c.Done()
//	go func() {
//		<-notify
//		// Remove this client from the map of attached clients
//		b.defunctClients <- messageChan
//	}()
//
//
//	for {
//		msg, open := <-messageChan
//		if !open {
//			// If our messageChan was closed, this means that
//			// the client has disconnected.
//			break
//		}
//
//		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
//		// Flush the response. This is only possible if the repsonse
//		// supports streaming.
//		f.Flush()
//		//err := sse.Encode(c.Writer, sse.Event{
//		//	Event: "message",
//		//	Data:  msg,
//		//})
//		//fmt.Println(err)
//	}
//
//	c.AbortWithStatus(http.StatusOK)
//}
