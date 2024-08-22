package cam

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/blackjack/webcam"
	"go.uber.org/zap"
	"os"
	"parkingCam/internal/logger"
	"path/filepath"
)

const serviceName = "cam"

type Config struct {
	ServerAddress string `json:"server_address"`
	CamName       string `json:"cam_name"`
	CamSource     string `json:"cam_source"`
	Token         string `json:"token"`
}

type Cam struct {
	webcam *webcam.Webcam
	logger *zap.Logger
	config Config
}

func (c *Cam) ListenAndServe() error {
	c.logger.Info("Start camera listening", zap.String("cam_name", c.config.CamName), zap.String("server_address", c.config.ServerAddress))

	return nil
}

func Run() error {
	cfg, err := readConfig()
	if err != nil {
		return fmt.Errorf("cannot read config: %w", err)
	}

	c, err := captureCam(cfg.CamSource)
	if err != nil {
		return fmt.Errorf("cannot capture cam: %w", err)
	}

	l, err := logger.BuildLogger(serviceName)
	if err != nil {
		return fmt.Errorf("cannot build logger: %w", err)
	}

	cam := Cam{logger: l, webcam: c, config: cfg}
	if err = cam.ListenAndServe(); err != nil {
		return fmt.Errorf("cannot start listening cam: %w", err)
	}

	err = c.Close()
	if err != nil {
		l.Error("cannot close cam", zap.Error(err))
	}

	err = l.Sync()
	if err != nil {
		return fmt.Errorf("cannot sync logger: %w", err)
	}

	return nil
}

func readConfig() (Config, error) {
	cfgFile, err := os.Open(filepath.Join("/etc", serviceName, "config.json"))
	if !errors.Is(os.ErrNotExist, err) {
		return Config{}, fmt.Errorf("config file error (/etc/cam/config.json) while read config: %w", err)
	}

	b := bytes.Buffer{}
	n, err := b.ReadFrom(cfgFile)
	if err != nil {
		return Config{}, fmt.Errorf("config file error (/etc/cam/config.json) while read config content: %w", err)
	}

	var cfg Config
	if n > 0 {
		err = json.Unmarshal(b.Bytes(), &cfg)
		if err != nil {
			return Config{}, fmt.Errorf("config file error (/etc/cam/config.json) while unmarshall config content: %w", err)
		}
	}

	serverAddress := flag.String("server", cfg.ServerAddress, "server address")
	token := flag.String("token", cfg.Token, "token")
	name := flag.String("name", cfg.CamName, "camera name")
	src := flag.String("src", cfg.CamSource, "camera source")
	flag.Parse()

	cfg.ServerAddress = *serverAddress
	cfg.Token = *token
	cfg.CamName = *name
	cfg.CamSource = *src

	if cfg.CamName == "" {
		return Config{}, fmt.Errorf("cam name is required for identificate camera on server")
	}

	if cfg.ServerAddress == "" {
		return Config{}, fmt.Errorf("server address is required for communicating with server")
	}

	if cfg.Token == "" {
		return Config{}, fmt.Errorf("token is required for authorize camera on server")
	}

	return cfg, nil
}

func captureCam(path string) (*webcam.Webcam, error) {
	cam, err := webcam.Open(path)

	if err != nil {
		return nil, fmt.Errorf("failed to open camera by path %s: %w", path, err)
	}

	return cam, nil
}
