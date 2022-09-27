package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/zx06/saut/internal/pkg/client/model"

	"github.com/gorilla/websocket"
	"github.com/zx06/saut/internal/pkg/assets_connector"
)

const (
	defaultRespDuration = time.Millisecond * 10
)

func WSHandler(c *websocket.Conn, a assets_connector.AssetsConnector) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.TODO())
	respTicker := time.NewTicker(defaultRespDuration)
	defer respTicker.Stop()
	err := a.Attach(ctx)
	if err != nil {
		log.Printf("attach assets_connector error: %s\n", err)
		cancel()
		return
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-respTicker.C:
				{
					var buf = make([]byte, 1024)
					n, err := a.Read(buf)
					if err != nil {
						log.Printf("read data err: %s \n", err)
						cancel()
						return
					}
					buf = buf[:n]
					data, err := model.NewWsTerminalOutputResponse(string(buf))
					if err != nil {
						log.Printf("marshal data err: %s \n", err)
						cancel()
						return
					}
					err = c.WriteJSON(data)
					if err != nil {
						log.Printf("write json err: %s \n", err)
						cancel()
						return
					}
				}
			case <-ctx.Done():
				//log.Println("stop read")
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				//log.Println("stop write")
				return
			default:
				{
					req := &model.WsTerminalRequest{}
					err := c.ReadJSON(req)
					if err != nil {
						log.Printf("read json err: %s \n", err)
						cancel()
						return
					}
					switch req.ReqType {
					case model.TerminalInput:
						{
							data, err := req.ParseTerminalInput()
							if err != nil {
								log.Printf("parse terminal input err: %s \n", err)
								cancel()
								return
							}
							_, err = a.Write([]byte(data))
							if err != nil {
								log.Printf("write data err: %s \n", err)
								cancel()
								return
							}
						}
					case model.TerminalResize:
						{
							data, err := req.ParseTerminalResize()
							if err != nil {
								log.Printf("parse terminal resize err: %s \n", err)
								cancel()
								return
							}
							err = a.WindowChange(data.H, data.W)
							if err != nil {
								log.Printf("window change err: %s \n", err)
								cancel()
								return
							}
						}
					}
				}
			}
		}
	}()
	wg.Wait()
	fmt.Println("-----")
}
