package Controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)
var bocker *SSEHandler
func init()  {
	log.SetFlags(log.Lshortfile | log.Ltime)
	bocker = NewSSEHandler()
}


type SSEHandler struct {
	clients map[chan string]bool

	newClients chan chan string

	defunctClients chan chan string

	messages chan string
}

func NewSSEHandler() *SSEHandler {
	b := &SSEHandler{
		clients:        make(map[chan string]bool),
		newClients:     make(chan (chan string)),
		defunctClients: make(chan (chan string)),
		messages:       make(chan string, 10),
	}

	return b
}

func (b *SSEHandler) HandleEvents() {
	go func() {
		for {
			select {
			case id := <-b.newClients:
				b.clients[id] = true
				log.Printf("Client added. %d registered clients", len(b.clients))
			case id := <-b.defunctClients:
				delete(b.clients, id)
				log.Printf("Removed client. %d registered clients", len(b.clients))
			case msg := <-b.messages:
				for id, _ := range b.clients {
					id <- msg
				}
			}
		}
	}()
}

func (b *SSEHandler) SendString(msg string) {
	b.messages <- msg
}

func (b *SSEHandler) SendJSON(obj interface{}) {
	tmp, err := json.Marshal(obj)
	if err != nil {
		log.Panic("Error while sending JSON object:", err)
	}
	b.messages <- string(tmp)
}

func (b *SSEHandler) Subscribe(c *gin.Context) {
	w := c.Writer
	f, ok := w.(http.Flusher)

	if !ok {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("Streaming unsupported"))
		return
	}

	messageChan := make(chan string)
	b.newClients <- messageChan
	go func() {
		b.messages <- "欢迎"+c.Query("name")+"进入了聊天室。"
	}()

	defer func() {
		b.defunctClients <- messageChan
	}()

	notify := w.CloseNotify()
	go func() {
		<-notify
		b.defunctClients <- messageChan
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin","*")

	for {
		msg, open := <-messageChan
		if !open {
			break
		}
		fmt.Fprintf(w, "data: %s\n\n", msg)
		f.Flush()
	}

	c.AbortWithStatus(http.StatusOK)
}

func Subscribe(c *gin.Context)  {
	bocker.HandleEvents()
	bocker.Subscribe(c)

	bocker.SendString("欢迎"+c.Query("name")+"进入了聊天室。")
}

func SendMsg(c *gin.Context)  {
	bocker.HandleEvents()
	bocker.SendString(c.Query("name")+":"+c.Query("message"))
}