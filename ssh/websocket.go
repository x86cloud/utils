package ssh

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"log"
	"net/http"
)

type Message struct {
	Op, Data   string
	Rows, Cols uint16
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (con *connection) SshClient(c *gin.Context) error {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	session, err := con.session()
	if err != nil {
		return err
	}
	defer session.Close()

	t := NewTTYHandler(ws, session)
	t.ws.SetCloseHandler(func(code int, text string) error {
		t.Close()
		return nil
	})

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err = session.RequestPty("xterm", 50, 180, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	t.ResizeEvent()

	if err = session.Shell(); err != nil {
		return err
	}

	if err = session.Wait(); err != nil {
		return err
	}
	return nil
}
