package base

import (
	"io/ioutil"
	"encoding/json"
)

//InitConfigFile 初始化配置文件信息
func InitJsonConfigFile(configFilePath string, out interface{}) error {
	content, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	err = FromByteJSON(content, out)
	if err != nil {
		return err
	}
	return nil
}

//ToJson JSON化字符串
func ToJson(obj interface{}) string {
	json, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(json)
}

//FromJSON 返回序列化为对象
func FromJSON(data string, obj interface{}) error {
	err := json.Unmarshal([]byte(data), obj)
	if err != nil {
		return err
	}
	return nil
}

//FromByteJSON 返回序列化对象
func FromByteJSON(data []byte, obj interface{}) error {
	err := json.Unmarshal(data, obj)
	if err != nil {
		return err
	}
	return nil
}