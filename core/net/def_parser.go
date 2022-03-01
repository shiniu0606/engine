package net

import (
	"reflect"
)

type ParserType int

const (
	ParserTypePB      ParserType = iota //protobuf
	ParserTypeJson                      //json
	ParserTypeMsgpack                   //msgpack
	ParserTypeRaw                       //不做任何解析
)

type ParseFunc func() interface{}

type IParser interface {
	GetType() ParserType
	UnPack(msg *Message) error
	Pack(v interface{}) []byte
	Register(cmd, act uint8, st interface{})
	RegisterAction(action string, st interface{})
}

type IParserFactory interface {
	Get() IParser
}

type Parser struct {
	ptype     ParserType
	typeMap   map[int]reflect.Type
	actionMap map[string]reflect.Type
}

func (r *Parser) GetType() ParserType {
	return r.ptype
}

//命令映射消息结构体
func (r *Parser) Register(cmd, act uint8, st interface{}) {
	rt := reflect.TypeOf(st).Elem()

	r.typeMap[CmdAct(cmd, act)] = rt
}

func (r *Parser) RegisterAction(action string, st interface{}) {
	rt := reflect.TypeOf(st).Elem()

	r.actionMap[action] = rt
}

func NewParser(Type ParserType) IParser {
	if Type == ParserTypePB {
		return &PBParser{
			Parser: Parser{
				ptype:   ParserTypePB,
				typeMap: make(map[int]reflect.Type)},
		}
	} else if Type == ParserTypeJson {
		return &JsonParser{
			Parser: Parser{
				ptype:   ParserTypeJson,
				typeMap: make(map[int]reflect.Type)},
		}
	}

	return nil
}
