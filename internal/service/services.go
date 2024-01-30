package service

import (
	"context"
	"device-worker-service/config"

	"go.uber.org/zap"
)

type IServer interface {
	StartServerListener()
	GetDevicesListener()
	GetPid(data string)
	GetDevices() error
	StopServer() error
}

type IDevice interface {
	TakeScreenshot(uui string) ([]byte, error)
}

type Services struct {
	ServerService IServer
	DeviceService IDevice
}

func NewServices(ctx context.Context, appConfig *config.Config, zapLogger *zap.Logger) *Services {

	services := &Services{
		ServerService: NewServerService(ctx, zapLogger, appConfig),
		DeviceService: NewDeviceService(ctx, zapLogger, appConfig),
	}

	go services.ServerService.StartServerListener()
	go services.ServerService.GetDevicesListener()

	return services
}
