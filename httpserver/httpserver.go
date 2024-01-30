package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func NewServer(handler *gin.Engine) *http.Server {

	return &http.Server{
		Handler:      handler,
		Addr:         ":8087",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

}
