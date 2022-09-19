package model

import "encoding/json"

type ReqType string

const (
	TerminalInput  ReqType = "TerminalInput"
	TerminalResize ReqType = "TerminalResize"
)

type RespType string

const (
	TerminalOutput RespType = "TerminalOutput"
)

type WsTerminalRequest struct {
	ReqType ReqType         `json:"req_type"`
	Data    json.RawMessage `json:"data"`
}

type WsTerminalRequestResize struct {
	H int `json:"h"`
	W int `json:"w"`
}

func (r WsTerminalRequest) ParseTerminalInput() (string, error) {
	var input string
	err := json.Unmarshal(r.Data, &input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func (r WsTerminalRequest) ParseTerminalResize() (*WsTerminalRequestResize, error) {
	var resize = &WsTerminalRequestResize{}
	err := json.Unmarshal(r.Data, resize)
	if err != nil {
		return nil, err
	}
	return resize, nil
}

type WsTerminalResponse struct {
	RespType RespType        `json:"resp_type"`
	Data     json.RawMessage `json:"data"`
}

func NewWsTerminalOutputResponse(data string) (*WsTerminalResponse, error) {
	var r = &WsTerminalResponse{
		RespType: TerminalOutput,
	}
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	r.Data = b
	return r, nil
}
