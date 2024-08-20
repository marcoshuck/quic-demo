package camera

import (
	"errors"
	vidio "github.com/AlexEidt/Vidio"
	"io"
)

type Camera interface {
	io.ReadCloser
}

type camera struct {
	adaptee *vidio.Camera
}

func (c *camera) Read(p []byte) (n int, err error) {
	if ok := c.adaptee.Read(); !ok {
		return 0, errors.New("failed to read camera video")
	}
	return copy(p, c.adaptee.FrameBuffer()), nil
}

func (c *camera) Close() error {
	c.adaptee.Close()
	return nil
}

func NewCamera(c *vidio.Camera) Camera {
	return &camera{
		adaptee: c,
	}
}
