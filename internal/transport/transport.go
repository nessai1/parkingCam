package transport

import (
	"fmt"
	srt "github.com/datarhei/gosrt"
	"go.uber.org/zap"
)

type Bucket struct {
	SourceID string
	Frame    []byte
}

type Sender struct {
}

type Receiver struct {
	logger     *zap.Logger
	passphrase string
}

func NewReceiver() *Receiver {
	return &Receiver{}
}

func (r *Receiver) Listen(f func(bucket Bucket)) error {
	ln, err := srt.Listen("srt", ":6000", srt.Config{
		Passphrase: r.passphrase,
	})

	if err != nil {
		return fmt.Errorf("cannot start listening transport port: %w", err)
	}

	for {
		req, err := ln.Accept2()
		if err != nil {
			r.logger.Error("cannot accept connection from camera", zap.Error(err))
			continue
		}

		go r.accept(req, f)
	}
}

func (r *Receiver) accept(request srt.ConnRequest, callback func(bucket Bucket)) {
	c, err := request.Accept()
	if err != nil {
		r.logger.Error("cannot accept connection from camera", zap.Error(err))
	}

	fmt.Println(c.StreamId())
	for {
		// TODO
	}
}
