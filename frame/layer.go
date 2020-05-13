package frame

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"io"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ILayer is an interface for a frame layer.
type ILayer interface {
	Write(frame *Frame)
	GetOutputChannel() <-chan Frame
}

// Layer contains a frame layer.
type Layer struct {
	transportLayer io.ReadWriter

	frameParser      *Parser
	parserInput      chan<- byte
	parserOutput     <-chan *ParseEvent
	acks, naks, cans <-chan bool

	l *zap.Logger

	pendingWrites chan *Frame
	frameOutput   chan Frame

	ctx context.Context
}

// NewFrameLayer will return a new frame layer.
func NewFrameLayer(ctx context.Context, transportLayer io.ReadWriter, logger *zap.Logger) (*Layer, error) {
	if _, ok := transportLayer.(io.ByteReader); !ok {
		return nil, errors.New("transport layer does not implement io.ByteReader")
	}

	parserInput := make(chan byte)
	parserOutput := make(chan *ParseEvent, 1)
	acks := make(chan bool, 1)
	naks := make(chan bool, 1)
	cans := make(chan bool, 1)

	parser := NewParser(ctx, parserInput, parserOutput, acks, naks, cans, logger)

	frameLayer := Layer{
		transportLayer: transportLayer,
		frameParser:    parser,
		parserInput:    parserInput,
		parserOutput:   parserOutput,
		acks:           acks,
		naks:           naks,
		cans:           cans,
		l:              logger,
		pendingWrites:  make(chan *Frame),
		frameOutput:    make(chan Frame, 5),
		ctx:            ctx,
	}

	go frameLayer.bgWork()
	go frameLayer.bgRead()

	return &frameLayer, nil

}

func (l *Layer) bgWork() {

	for {
		select {
		case frameIn := <-l.parserOutput:
			l.l.Debug("parser output received")

			if frameIn.status == ParseOk {
				l.sendAck()
				l.l.Debug("received frame successfully, writing output", zap.String("frame", spew.Sdump(frameIn)))
				l.frameOutput <- frameIn.frame
			} else if frameIn.status == ParseNotOk {
				l.l.Warn("received frame, parse not ok", zap.String("frame", spew.Sdump(frameIn)))
				l.sendNak()
			} else {
				// @todo handle timeout(?)
			}

		case <-l.acks:
			l.l.Debug("rx ack")
		case <-l.naks:
			l.l.Debug("rx nak")
		case <-l.cans:
			l.l.Debug("rx can")

		case frameToWrite := <-l.pendingWrites:
			l.l.Debug("frame received, writing to transport")
			// this method never returns an error, so ignore it
			buf, _ := frameToWrite.MarshalBinary()
			attempts := time.Duration(0)
		retry:
			l.l.Debug("writing to transport", zap.Int64("attempt", int64(attempts)))

			l.writeToTransport(buf)
			// TODO: this needs to time out

		drain:
			for {
				select {
				case <-l.acks:
					l.l.Debug("received ack")
					attempts = 100 // break out of retry loop
					break drain
				case <-l.cans:
					l.l.Warn("received can")
				case <-l.naks:
					l.l.Warn("received nak")
				case frameIn := <-l.parserOutput:
					l.l.Warn("received frame while waiting for ack", zap.String("frame", spew.Sdump(frameIn)))
					l.sendNak()

					// TODO wait 1500 ms for ack
					// in case of fail, resent after 100ms + att*1s

				case <-time.After(100 * time.Millisecond + attempts * time.Second):
					l.l.Error("ack timed out")
					break drain
				}
			}
			attempts++
			if attempts <= 3 {
				goto retry
			}
			// how to report error up?

		case <-l.ctx.Done():
			l.l.Info("closing frame layer bg work")
			return
		}
	}
}

func (l *Layer) Write(frame *Frame) {
	go func() {
		l.pendingWrites <- frame
	}()
}

// GetOutputChannel will return the output channel.
func (l *Layer) GetOutputChannel() <-chan Frame {
	return l.frameOutput
}

func (l *Layer) bgRead() {
	for {
		byt, err := l.transportLayer.(io.ByteReader).ReadByte()
		if err == io.EOF {
			// TODO: handle EOF
			return
		} else if err != nil {
			// TODO: handle more gracefully
			l.l.Fatal("error reading from transport", zap.String("err", err.Error()))
		}
		//l.l.Debug("got byte", zap.Int("byte", int(byt)))

		l.parserInput <- byt
	}
}

func (l *Layer) writeToTransport(buf []byte) (int, error) {
	l.l.Debug("send data", zap.String("buf", spew.Sdump(buf)))
	return l.transportLayer.Write(buf)
}

func (l *Layer) sendAck() error {
	l.l.Debug("send ack")
	_, err := l.transportLayer.Write([]byte{HeaderAck})
	return err
}

func (l *Layer) sendNak() error {
	l.l.Debug("send nak")
	_, err := l.transportLayer.Write([]byte{HeaderNak})
	return err
}
