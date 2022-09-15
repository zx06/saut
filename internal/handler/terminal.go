package handler

import (
	"log"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type DataType string

const (
	TerminalInput  DataType = "TerminalInput"
	TerminalOutput DataType = "TerminalOutput"
)

type RequestWS struct {
	Type DataType `json:"type"`
	Data []byte   `json:"data"`
}

type ResponseWS struct {
	Type DataType `json:"type"`
	Data []byte   `json:"data"`
}

func TerminalHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err,
		})
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"message": err,
			})
		}
		return
	}(conn)
	terminalHandlerWS(conn)

}

func terminalHandlerWS(c *websocket.Conn) {
	respTicker := time.NewTicker(time.Millisecond * 20)
	defer respTicker.Stop()
	sshCfg := &ssh.ClientConfig{
		User: "****",
		Auth: []ssh.AuthMethod{
			ssh.Password("****"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConn, err := ssh.Dial("tcp", "****", sshCfg)
	if err != nil {
		log.Println(err)
		return
	}
	defer sshConn.Close()
	sess, err := sshConn.NewSession()
	if err != nil {
		log.Println(err)
		return
	}
	defer sess.Close()
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}
	if err = sess.RequestPty("linux", 32, 160, modes); err != nil {
		log.Println("request pty error: %s", err.Error())
		return
	}
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Println(err)
		return
	}
	stdout, err := sess.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	//stderr, err := sess.StderrPipe()
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	if err := sess.Shell(); err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			<-respTicker.C
			var buf = make([]byte, 1024)
			_, _ = stdout.Read(buf)

			err := c.WriteJSON(&ResponseWS{
				Type: TerminalOutput,
				Data: buf,
			})
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	for {
		var (
			req RequestWS
		)
		err := c.ReadJSON(&req)
		if err != nil {
			log.Println(err)
			break
		}
		_, err = stdin.Write(req.Data)
		if err != nil {
			log.Println(err)
			return
		}
	}
	if err := sess.Wait(); err != nil {
		log.Println(err)
		return
	}
}
