package Controllers

import (
	"fmt"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"learn-golang/Services"
	"net/http"
	"time"
)

func TestInsert(c *gin.Context) {


	fmt.Println(c)
	var testService Services.Test

	err := c.ShouldBindJSON(&testService)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := testService.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"message": "Insert() error!",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"message": "success",
		"data": id,
	})

}

func Sse(c *gin.Context) {
	fmt.Println(sse.ContentType)
	c.Writer.Header().Set("Content-Type","text/event-stream")
	c.Writer.Header().Set("Cache-Control","no-cache")
	c.Writer.Header().Set("Connection","keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin","*")
	for {
		// data can be a primitive like a string, an integer or a float
		sse.Encode(c.Writer, sse.Event{
			Event: "message",
			Data:  "some data\nmore data",
		})
		time.Sleep(1000)
	}


	// also a complex type, like a map, a struct or a slice
	//sse.Encode(c.Writer, sse.Event{
	//	Id:    "124",
	//	Event: "message",
	//	Data: map[string]interface{}{
	//		"user":    "manu",
	//		"date":    time.Now().Unix(),
	//		"content": "hi!",
	//	},
	//})
}