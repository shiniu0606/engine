package net

import (
	"errors"
	"reflect"

	"github.com/golang/protobuf/proto"

	base "github.com/shiniu0606/engine/core/base"
)

var ErrMsgPbUnPack = errors.New("message pb unpack error")

type PBParser struct {
	Parser
}

func (r *PBParser) UnPack(msg *Message) error {
	if msg == nil {
		return ErrMsgPbUnPack
	}

	if msg.Head == nil {
		return ErrMsgPbUnPack
	}

	if p, ok := r.typeMap[msg.Head.CmdAct()]; ok {
		st := reflect.New(p).Interface()
		if st != nil {
			if len(msg.Data) > 0 {
				err := PBUnPack(msg.Data, st)
				msg.UserData = st
				if err != nil {
					base.LogInfo("PBUnPack error:%v", err)
					return ErrMsgPbUnPack
				}
			}
			return nil
		}
	}

	return nil
}

func (r *PBParser) Pack(v interface{}) []byte {
	data, _ := PBPack(v)
	return data
}

func PBUnPack(data []byte, msg interface{}) error {
	if data == nil || msg == nil {
		return nil
	}

	err := proto.Unmarshal(data, msg.(proto.Message))
	if err != nil {
		return err
	}
	return nil
}

func PBPack(msg interface{}) ([]byte, error) {
	if msg == nil {
		return nil, nil
	}

	data, err := proto.Marshal(msg.(proto.Message))

	return data, err
}
