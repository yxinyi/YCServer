package YDecode

import (
	"github.com/json-iterator/go"
)

const (
	DECODE_TYPE_JSON uint32 = 0
)

type Inter interface {
	Marshal(interface{})  ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var g_decode_pool = make(map[uint32]Inter)

func init() {
	g_decode_pool[DECODE_TYPE_JSON] = &Json{}
}

type Json struct{}

func (j *Json) Marshal(val interface{}) ([]byte, error) {
	return jsoniter.Marshal(val)
}
func (j *Json) Unmarshal(data_ []byte, val_ interface{}) error {
	return jsoniter.Unmarshal(data_, val_)
}

func Marshal(type_ uint32,val interface{} )([]byte, error){
	return g_decode_pool[type_].Marshal(val)
}

func Unmarshal(type_ uint32,data_ []byte, val_ interface{})error{
	return g_decode_pool[type_].Unmarshal(data_,val_)
}