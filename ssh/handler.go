package ssh

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"k8s.io/klog/v2"
)
type Message struct {
	Op, Data   string
	Rows, Cols uint16
}

type TTYReadWriter interface {
	Read(m *Message) error
	Writer(m *Message) error
}

type TTYHandler struct {
	rw TTYReadWriter
	resize  chan struct{ Cols, Rows uint16 }
	CloseCh chan struct{}
}

func NewTTYHandler(rw TTYReadWriter) *TTYHandler {
	tty := &TTYHandler{
		rw:      rw,
		CloseCh: make(chan struct{}),
		resize:  make(chan struct{ Cols, Rows uint16 }),
	}
	return tty
}

func (t TTYHandler) ResizeEvent(session *ssh.Session) {
	go func() {
		for true {
			select {
			case resize := <-t.resize:
				if err := session.WindowChange(int(resize.Rows), int(resize.Cols)); err != nil {
					if err.Error() == "EOF" {
						return
					}
					klog.Warning("ssh resize failed, err: ", err)
				}
			case <-t.CloseCh:
				return
			}
		}
	}()
}

func (t TTYHandler) Read(p []byte) (int, error) {
	msg := &Message{}
	err := t.rw.Read(msg)
	if err != nil {
		return 0, err
	}

	switch msg.Op {
	case "stdin":
		return copy(p, msg.Data), nil
	case "resize":
		t.resize <- struct{ Cols, Rows uint16 }{Cols: msg.Cols, Rows: msg.Rows}
	default:
		return 0, errors.New("unKnown type: " + msg.Op)
	}
	return 0, nil
}

func (t TTYHandler) Write(p []byte) (int, error) {
	if err := t.rw.Writer(&Message{
		Op:   "stdout",
		Data: string(p),
	}); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (t TTYHandler) Close() {
	close(t.CloseCh)
	close(t.resize)
}
