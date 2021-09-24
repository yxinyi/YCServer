package YNet

import (
	"encoding/binary"
	"github.com/json-iterator/go"
	"io"
	"reflect"
	"strings"
	"unsafe"
)

const (
	MESSAGE_TYPE_JSON  = 1
	MESSAGE_TYPE_PROTO = 2
)

type NetMsgPack struct {
	M_msg_name string
	M_msg_data []byte
}

const (
	const_type_size   = int(unsafe.Sizeof(uint32(0)))
	const_length_size = int(unsafe.Sizeof(uint32(0)))
)

func NewNetMsgPack() *NetMsgPack {
	return &NetMsgPack{
	}
}
func NewNetMsgPackWithJson(json_ interface{}) *NetMsgPack {
	_msg := NewNetMsgPack()
	
	_msg_name := reflect.TypeOf(json_).String()
	_split_idx := strings.Index(_msg_name,".")
	_msg_name = _msg_name[_split_idx+1:]
	_msg.M_msg_name = _msg_name
	_byte, _err := jsoniter.Marshal(json_)
	if _err != nil {
		return nil
	}
	_msg.M_msg_data = _byte
	return _msg
}

func (pack *NetMsgPack) ToByteStream() []byte {
	_name_length := len(pack.M_msg_name)
	_data_length := len(pack.M_msg_data)
	
	_total_length := const_type_size + const_length_size + _name_length + _data_length
	_stream_byte := make([]byte, _total_length)
	
	binary.LittleEndian.PutUint32(_stream_byte, uint32(_name_length))
	binary.LittleEndian.PutUint32(_stream_byte[const_type_size:], uint32(_data_length))
	copy(_stream_byte[uint32(const_type_size+const_length_size):], pack.M_msg_name[:])
	copy(_stream_byte[uint32(const_type_size+const_length_size+_name_length):], pack.M_msg_data[:])
	return _stream_byte
}

func (pack *NetMsgPack) InitFromIO(io_ io.Reader) bool {
	_type_byte := make([]byte, const_type_size+const_length_size)
	
	_len, err := io.ReadFull(io_, _type_byte)
	if _len == 0 {
		return false
	}
	if err != nil {
		return false
	}
	
	_name_length := binary.LittleEndian.Uint32(_type_byte[0:const_type_size])
	_data_length := binary.LittleEndian.Uint32(_type_byte[const_type_size : const_type_size+const_type_size])
	
	_msg_name_byte := make([]byte, _name_length)
	_len, err = io.ReadFull(io_, _msg_name_byte)
	if err != nil {
		return false
	}
	pack.M_msg_name = string(_msg_name_byte)
	
	pack.M_msg_data = make([]byte, _data_length)
	_len, err = io.ReadFull(io_, pack.M_msg_data)
	if err != nil {
		return false
	}
	
	return true
}
