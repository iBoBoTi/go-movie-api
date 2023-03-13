package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Data struct {
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Errors    string      `json:"errors,omitempty"`
	Status    string      `json:"status"`
}

func JSON(c *gin.Context, status int, message string, data interface{}, err error) {
	errMessage := ""
	if err != nil {
		errMessage = err.Error()
	}
	responsedata := Data{
		Message:   message,
		Data:      data,
		Errors:    errMessage,
		Status:    http.StatusText(status),
		Timestamp: time.Now().Format(time.RFC850),
	}

	c.JSON(status, responsedata)
}
