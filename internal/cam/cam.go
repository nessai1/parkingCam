package cam

import (
	"fmt"
	"github.com/blackjack/webcam"
)

func Run() error {
	cam, err := captureCam()
	if err != nil {
		return fmt.Errorf("cannot capture cam: %w", err)
	}

	defer cam.Close()

	err = cam.StartStreaming()
	if err != nil {
		return fmt.Errorf("failed to start streaming: %v", err)
	}

	for {
		err = cam.WaitForFrame(1)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			continue
		default:
			return fmt.Errorf("failed to wait for frame: %v", err)
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			return fmt.Errorf("failed to read frame: %v", err)
		}

		fmt.Println(len(frame))
	}
}

func captureCam() (*webcam.Webcam, error) {
	// todo: choose cam dialog
	camPath := "/dev/video0"
	cam, err := webcam.Open(camPath)

	if err != nil {
		return nil, fmt.Errorf("failed to open camera by path %s: %w", camPath, err)
	}

	return cam, nil
}
