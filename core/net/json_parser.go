package net

import (
	"encoding/json"
	"errors"
	"reflect"

	base "github.com/shiniu0606/engine/core/base"
)

var ErrMsgJsonUnPack = errors.New("message json unpack error")

type JsonParser struct {
	Parser
}

type JsonMessage struct {
	Action string `json:"action"`
}

func (r *JsonParser) UnPack(msg *Message) error {
	if msg == nil {
		return ErrMsgJsonUnPack
	}

	if msg.Head == nil {
		var actionMessage JsonMessage
		err := json.Unmarshal(msg.Data, &actionMessage)
		if err != nil {
			base.LogInfo("JsonUnPack error:%v", err)
			return ErrMsgJsonUnPack
		}
		if p, ok := r.actionMap[actionMessage.Action]; ok {
			st := reflect.New(p).Interface()
			if st != nil {
				if len(msg.Data) > 0 {
					err := JsonUnPack(msg.Data, st)
					msg.UserData = st
					if err != nil {
						base.LogInfo("JsonUnPack error:%v", err)
						return ErrMsgJsonUnPack
					}
				}
				return nil
			}
		}
		return ErrMsgJsonUnPack
	}

	if p, ok := r.typeMap[msg.Head.CmdAct()]; ok {
		st := reflect.New(p).Interface()
		if st != nil {
			if len(msg.Data) > 0 {
				err := JsonUnPack(msg.Data, st)
				msg.UserData = st
				if err != nil {
					base.LogInfo("JsonUnPack error:%v", err)
					return ErrMsgJsonUnPack
				}
			}
			return nil
		}
	}

	return nil
}

func (r *JsonParser) Pack(v interface{}) []byte {
	data, _ := JsonPack(v)
	return data
}

func JsonUnPack(data []byte, msg interface{}) error {
	if data == nil || msg == nil {
		return nil
	}

	err := json.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	return nil
}

func JsonPack(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}

	data, err := json.Marshal(msg)

	return data, err
}
