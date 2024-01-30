package handler

import (
	_ "device-worker-service/docs"
	"device-worker-service/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type Handler struct {
	services  *service.Services
	zapLogger *zap.Logger
}

func NewHandler(services *service.Services, zapLogger *zap.Logger) *Handler {

	return &Handler{
		services:  services,
		zapLogger: zapLogger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	r := gin.Default()
	r.Use(gin.Recovery())

	r.GET("/device-worker/ping", h.Pong)
	r.GET("/start-server", h.StartServer)
	r.GET("/stop-server", h.StopServer)
	r.GET("/screenshot/:uuid", h.ScreenShot)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// @Summary Ping Pong
// @Description Get files by providing the filename as a query parameter.
// @Success 200 {json}  "Successful response"
func (h *Handler) Pong(c *gin.Context) {

	c.JSON(http.StatusOK, "pong")
}

func (h *Handler) StartServer(c *gin.Context) {

	if err := h.services.ServerService.GetDevices(); err != nil {
		h.zapLogger.Error("error in start server", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in start server"})
		return
	}

	c.JSON(http.StatusOK, "started server")
}

func (h *Handler) StopServer(c *gin.Context) {

	c.JSON(http.StatusOK, "pong")
}

// @Summary Get files by filename
// @Description Get files by providing the filename as a query parameter.
// @Produce png
// @Param filename query string true "File name" default("example.png")
// @Success 200 {string} image/png "Successful response"
// @Failure 400 {json}  "file name not specified"
// @Failure 404 {json}  "error file not found"
// @Failure 500 {json}  "Internal Server Error" example({"error":"Internal Server Error"})
// @Router /screenshot/:uuid [get]
func (h *Handler) ScreenShot(c *gin.Context) {

	uuid := c.Param("uuid")

	screenShotBytes, err := h.services.DeviceService.TakeScreenshot(uuid)
	if err != nil {
		h.zapLogger.Error("error in take screenshot service", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "error in take screenshot service"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Length", fmt.Sprint(len(screenShotBytes)))

	c.Data(http.StatusOK, "image/png", screenShotBytes)
}
