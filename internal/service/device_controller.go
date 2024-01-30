package service

import (
	"bytes"
	"context"
	"device-worker-service/config"
	"device-worker-service/traits"
	"fmt"
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

var (
	outputPath = "./screens/"
)

type DeviceService struct {
	ctx       context.Context
	zapLogger *zap.Logger
	appConfig *config.Config
}

func NewDeviceService(ctx context.Context, zapLogger *zap.Logger, appConfig *config.Config) *DeviceService {

	deviceService := &DeviceService{
		ctx:       ctx,
		zapLogger: zapLogger,
		appConfig: appConfig,
	}

	return deviceService
}

func (s *DeviceService) TakeScreenshot(uuid string) ([]byte, error) {

	dataJson, err := traits.ParseUrl(s.appConfig.LocalUrl)
	if err != nil {
		return nil, err
	}

	res, err := traits.ParseDevice(dataJson)
	if err != nil {
		return nil, err
	}

	rsdHost, rsdPort := traits.FindDataByUUId(res, uuid)

	fileName := outputPath + uuid + time.Now().String() + "_screenshot.png"

	cmd := exec.Command("pymobiledevice3", "developer", "dvt", "screenshot", fileName, "--rsd", rsdHost, rsdPort)

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running pymobiledevice3: %v, stderr: %s", err, stderrBuf.String())
	}

	screenshotBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("error reading screenshot file: %v", err)
	}

	return screenshotBytes, nil
}
