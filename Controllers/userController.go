package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"learn-golang/Services"
	"net/http"
)

func UserInsert(c *gin.Context) {
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
