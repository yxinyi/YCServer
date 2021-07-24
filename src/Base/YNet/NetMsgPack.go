package YNet

import (
	"encoding/binary"
	"io"
	"unsafe"
)

const (
	MESSAGE_TYPE_JSON  = 1
	MESSAGE_TYPE_PROTO = 2
)

type NetMsgPack struct {
	M_msg_id     uint32
	M_msg_length uint32
	M_msg_data   []byte
}

const (
	const_type_size   = int(unsafe.Sizeof(uint32(0)))
	const_length_size = int(unsafe.Sizeof(uint32(0)))
)

func NewMessagePack() *NetMsgPack {
	return &NetMsgPack{
	}
}

func (pack *NetMsgPack) ToByteStream() []byte {
	_total_length := const_type_size + const_length_size + len(pack.M_msg_data)
	_stream_byte := make([]byte, _total_length)

	binary.LittleEndian.PutUint32(_stream_byte, pack.M_msg_id)
	binary.LittleEndian.PutUint32(_stream_byte[const_type_size:], pack.M_msg_length)
	copy(_stream_byte[uint32(const_type_size + const_length_size):], pack.M_msg_data[:])
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

	pack.M_msg_id = binary.LittleEndian.Uint32(_type_byte[0:const_type_size])
	pack.M_msg_length = binary.LittleEndian.Uint32(_type_byte[const_type_size:const_type_size+const_length_size])

	pack.M_msg_data = make([]byte, pack.M_msg_length)
	_len, err = io.ReadFull(io_, pack.M_msg_data)
	if err != nil {
		return false
	}
	return true
}
