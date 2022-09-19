package assets_connector

import (
	"context"
	"fmt"
	"io"
	"log"

	"golang.org/x/crypto/ssh"
)

var _ AssetsConnector = (*SSHConnector)(nil)

var defaultTerminalModes = ssh.TerminalModes{
	ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
	ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
	ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
}

type SSHConnector struct {
	sshClient *ssh.Client
	session   *ssh.Session
	stdin     io.Writer
	stdout    io.Reader
}

func NewSSHConnector(sshCfg *ssh.ClientConfig, sshAddr string) (AssetsConnector, error) {
	uc := &SSHConnector{}
	sshClient, err := ssh.Dial("tcp", sshAddr, sshCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ssh :%w", err)
	}
	uc.sshClient = sshClient
	sess, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("new session error :%w", err)
	}
	uc.session = sess
	//uc.stdin, err = sess.StdinPipe()
	//if err != nil {
	//	return nil, fmt.Errorf("set stdin error :%w", err)
	//}
	//uc.stdout, err = sess.StdoutPipe()
	//if err != nil {
	//	return nil, fmt.Errorf("set stdout error :%w", err)
	//}
	return uc, nil
}

func (c *SSHConnector) Read(p []byte) (n int, err error) {
	return c.stdout.Read(p)
}

func (c *SSHConnector) Write(p []byte) (n int, err error) {
	return c.stdin.Write(p)
}

func (c *SSHConnector) Close() error {
	err := c.session.Close()
	if err != nil {
		return err
	}
	err = c.sshClient.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *SSHConnector) Attach(ctx context.Context) error {
	err := c.session.RequestPty("linux", 40, 80, defaultTerminalModes)
	if err != nil {
		return err
	}
	c.stdin, err = c.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("set stdin error :%w", err)
	}
	c.stdout, err = c.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("set stdout error :%w", err)
	}
	err = c.session.Shell()
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		c.Close()
	}()
	go func() {
		err := c.session.Wait()
		if err != nil {
			log.Printf("wait session error: %s\n", err)
			c.Close()
		}
	}()
	return nil
}

func (c *SSHConnector) WindowChange(h int, w int) error {
	return c.session.WindowChange(h, w)
}
