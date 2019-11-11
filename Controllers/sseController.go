// Server-side events handler for Gin.
// Based on work by Kyle L. Jensen
// Source: https://github.com/kljensen/golang-html5-sse-example

package Controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var bocker = NewSSEHandler()

type SSEHandler struct {
	clients map[chan string]bool

	newClients chan chan string

	defunctClients chan chan string

	messages chan string
}

// Make a new SSEHandler instance.
func NewSSEHandler() *SSEHandler {
	b := &SSEHandler{
		clients:        make(map[chan string]bool),
		newClients:     make(chan (chan string)),
		defunctClients: make(chan (chan string)),
		messages:       make(chan string, 10), // buffer 10 msgs and don't block sends
	}
	return b
}

// Start handling new and disconnected clients, as well as sending messages to
// all connected clients.
func (b *SSEHandler) HandleEvents(id string) {
	go func() {
		for {
			select {
			case id := <-b.newClients:
				b.clients[id] = true
				log.Printf("Client added. %d registered clients", len(b.clients))
			case id := <-b.defunctClients:

				delete(b.clients, id)
				close(id)
				log.Printf("Removed client. %d registered clients", len(b.clients))
			case msg := <-b.messages:
				log.Println(b.clients)
				for id, _ := range b.clients {
					log.Println(id)
					id <- msg

				}
			}
		}
	}()
}

// Send out a simple string to all clients.
func (b *SSEHandler) SendString(msg string) {
	b.messages <- msg
}

// Send out a JSON string object to all clients.
func (b *SSEHandler) SendJSON(obj interface{}) {
	tmp, err := json.Marshal(obj)
	if err != nil {
		log.Panic("Error while sending JSON object:", err)
	}
	b.messages <- string(tmp)
}

// Subscribe a new client and start sending out messages to it.
func (b *SSEHandler) Subscribe(c *gin.Context) {
	w := c.Writer
	f, ok := w.(http.Flusher)

	if !ok {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Streaming unsupported"))
		return
	}

	messageChan := make(chan string)
	b.newClients <- messageChan
	notify := c.Done()
	go func() {
		<-notify
		b.defunctClients <- messageChan
	}()
log.Println(b.clients)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin","*")

	for {
		msg, open := <-messageChan
		if !open {
			break
		}

		fmt.Fprintf(w, "data: Message: %s\n\n", msg)
		f.Flush()
	}

	c.AbortWithStatus(http.StatusOK)
}

func Subscribe(c *gin.Context)  {
	log.Println(c.Query("id"))
	bocker.HandleEvents(c.Query("id"))
	go func() {
		bocker.Subscribe(c)
	}()


}

func SendMsg(c *gin.Context)  {
	bocker.HandleEvents(c.Query("id"))
	go func() {
		bocker.SendString("hello world."+c.Query("id"))
	}()
}