package service

import (
	"bytes"
	"context"
	"device-worker-service/config"
	"device-worker-service/traits"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"go.uber.org/zap"
)

type ServerService struct {
	ctx       context.Context
	zapLogger *zap.Logger
	appConfig *config.Config
	Pid       int
	dataCmd   chan string
}

func NewServerService(ctx context.Context, zapLogger *zap.Logger, appConfig *config.Config) *ServerService {

	serverService := &ServerService{
		ctx:       ctx,
		zapLogger: zapLogger,
		appConfig: appConfig,
		dataCmd:   make(chan string, 1),
	}

	return serverService
}

func (s *ServerService) GetDevicesListener() {

	select {
	case data := <-s.dataCmd:
		s.GetPid(data)
	case <-time.After(5 * time.Second):
		s.GetDevices()
	}

}

func (s *ServerService) StartServerListener() {

	s.zapLogger.Info("started start-server service")

	cmd := exec.Command("sudo", "python3", "-m", "pymobiledevice3", "remote", "tunneld")

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	fmt.Fprintln(os.Stderr, "Command executed successfully.")

	if err := cmd.Start(); err != nil {
		s.zapLogger.Error("error starting command", zap.Error(err))
		return
	}

	if err := cmd.Wait(); err != nil {
		s.zapLogger.Error("error waiting for command", zap.Error(err))
		return
	}

	if stderrBuf.Len() > 0 {
		s.zapLogger.Error("error in command execution", zap.String("stderr", stderrBuf.String()))
		return
	}

	output, err := ioutil.ReadAll(&stdoutBuf)
	if err != nil {
		s.zapLogger.Error("error reading command output", zap.Error(err))
		return
	}

	s.dataCmd <- string(output)
}

func (s *ServerService) GetPid(data string) {

	pid, err := traits.ExtractPID(string(data))
	if err != nil {
		s.zapLogger.Error("error extracting PID from log", zap.Error(err))
		return
	}

	s.Pid = pid
}

func (s *ServerService) GetDevices() error {

	resp, err := traits.ParseUrl(s.appConfig.LocalUrl)
	if err != nil {
		return err
	}

	devices, err := traits.ParseDevice(resp)
	if err != nil {
		return err
	}

	fmt.Println(devices)

	return nil
}

func (s *ServerService) StopServer() error {

	return nil
}
